package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "koushoku/config"

	"koushoku/cache"
	"koushoku/database"
	"koushoku/errs"
	"koushoku/models"
	"koushoku/modext"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	ArchiveCols = models.ArchiveColumns
	ArchiveRels = models.ArchiveRels
)

func insertArchive(archive *modext.Archive, upsert bool) (*modext.Archive, error) {
	if archive == nil {
		return nil, nil
	} else if len(archive.Path) == 0 {
		return nil, errs.ArchivePathRequired
	}

	selectMods := []QueryMod{
		Where("archive.slug ILIKE ? AND archive.expunged IS FALSE", archive.Slug),
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Circles),
		Load(ArchiveRels.Magazines),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	}

	var q []string
	var args []any

	if len(archive.Artists) > 0 {
		for _, artist := range archive.Artists {
			q = append(q, rawSqlArtistsMatch)
			args = append(args, Slugify(artist.Name))
		}
		selectMods = append(selectMods,
			Where(JoinOR(q...), args...))
	} else if len(archive.Magazines) > 0 {
		for _, magazine := range archive.Magazines {
			q = append(q, rawSqlMagazinesMatch)
			args = append(args, Slugify(magazine.Name))
		}
		selectMods = append(selectMods,
			Where(JoinOR(q...), args...))
	} else if len(archive.Circles) > 0 {
		for _, circle := range archive.Circles {
			q = append(q, rawSqlCirclesMatch)
			args = append(args, Slugify(circle.Name))
		}
		selectMods = append(selectMods,
			Where(JoinOR(q...), args...))
	} else {
		selectMods = append(selectMods,
			Where("archive.path = ?", archive.Path))
	}

	model, err := models.Archives(selectMods...).OneG()
	if upsert {
		if err != nil && err != sql.ErrNoRows {
			return nil, errs.Unknown
		}
	} else if err == nil {
		return modext.NewArchive(model), nil
	} else {
		return nil, errs.ArchiveNotFound
	}

	tx, err := database.Conn.Begin()
	if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}

	var isDuplicate bool
	if model == nil {
		model = &models.Archive{Title: archive.Title, Slug: archive.Slug}
		if archive.CreatedAt > 0 {
			model.CreatedAt = time.Unix(archive.CreatedAt, 0)
			model.UpdatedAt = model.CreatedAt
		}
	} else {
		isDuplicate = true
		model.UpdatedAt = time.Unix(archive.CreatedAt, 0)
	}

	model.Path = archive.Path
	model.Pages = archive.Pages
	model.Size = archive.Size

	op := model.Insert
	if isDuplicate {
		op = model.Update
	}

	err = op(tx, boil.Infer())
	if err == nil {
		err = PopulateArchiveRels(tx, model, archive)
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}

	if err != nil {
		log.Println(errors.WithStack(err))
		return nil, errs.Unknown
	}

	// TODO: Purge cache
	return modext.NewArchive(model), nil
}

func CreateArchive(archive *modext.Archive) (*modext.Archive, error) {
	return insertArchive(archive, true)
}

func UpdateArchive(archive *modext.Archive) (*modext.Archive, error) {
	return insertArchive(archive, false)
}

type GetArchiveOptions struct {
	Preloads []string `form:"preload" json:"1,omitempty"`
}

type GetArchiveResult struct {
	Archive *modext.Archive `json:"archive,omitempty"`
	Err     error           `json:"error,omitempty"`
}

func GetArchive(id int64, opts GetArchiveOptions) (result *GetArchiveResult) {
	opts.Preloads = validateArchiveRels(opts.Preloads)

	cacheKey := makeCacheKey(opts)
	if c, err := cache.Archives.GetWithPrefix(id, cacheKey); err == nil {
		return c.(*GetArchiveResult)
	}

	result = &GetArchiveResult{}
	defer func() {
		if result.Archive != nil || result.Err != nil {
			cache.Archives.RemoveWithPrefix(id, cacheKey)
			cache.Archives.SetWithPrefix(id, cacheKey, result, 0)
		}
	}()

	selectQueries := []QueryMod{Where("id = ?", id), And("published_at IS NOT NULL")}
	for _, v := range opts.Preloads {
		if v == ArchiveRels.Artists || v == ArchiveRels.Circles || v == ArchiveRels.Tags {
			selectQueries = append(selectQueries, Load(v, OrderBy("name ASC")))
		} else {
			selectQueries = append(selectQueries, Load(v))
		}
	}

	archive, err := models.Archives(selectQueries...).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			result.Err = errs.ArchiveNotFound
			return
		}
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Archive = modext.NewArchive(archive).LoadRels(archive)
	return
}

type GetArchivesOptions struct {
	Path string `json:"0,omitempty"`

	TitleMatch    string `json:"1,omitempty"`
	TitleWildcard string `json:"2,omitempty"`

	ArtistsMatch            []string `json:"3,omitempty"`
	ArtistsMatchAnd         []string `json:"4,omitempty"`
	ArtistsWildcard         []string `json:"5,omitempty"`
	ArtistsWildcardAnd      []string `json:"6,omitempty"`
	ExcludedArtistsMatch    []string `json:"7,omitempty"`
	ExcludedArtistsWildcard []string `json:"8,omitempty"`

	CirclesMatch            []string `json:"9,omitempty"`
	CirclesMatchAnd         []string `json:"10,omitempty"`
	CirclesWildcard         []string `json:"11,omitempty"`
	CirclesWildcardAnd      []string `json:"12,omitempty"`
	ExcludedCirclesMatch    []string `json:"13,omitempty"`
	ExcludedCirclesWildcard []string `json:"14,omitempty"`

	MagazinesMatch            []string `json:"15,omitempty"`
	MagazinesMatchAnd         []string `json:"16,omitempty"`
	MagazinesWildcard         []string `json:"17,omitempty"`
	MagazinesWildcardAnd      []string `json:"18,omitempty"`
	ExcludedMagazinesMatch    []string `json:"19,omitempty"`
	ExcludedMagazinesWildcard []string `json:"20,omitempty"`

	ParodiesMatch            []string `json:"21,omitempty"`
	ParodiesMatchAnd         []string `json:"22,omitempty"`
	ParodiesWildcard         []string `json:"23,omitempty"`
	ParodiesWildcardAnd      []string `json:"24,omitempty"`
	ExcludedParodiesMatch    []string `json:"25,omitempty"`
	ExcludedParodiesWildcard []string `json:"26,omitempty"`

	TagsMatch            []string `json:"27,omitempty"`
	TagsMatchAnd         []string `json:"28,omitempty"`
	TagsWildcard         []string `json:"29,omitempty"`
	TagsWildcardAnd      []string `json:"30,omitempty"`
	ExcludedTagsMatch    []string `json:"31,omitempty"`
	ExcludedTagsWildcard []string `json:"32,omitempty"`

	PagesEq  int `json:"33,omitempty"`
	PagesGt  int `json:"34,omitempty"`
	PagesGte int `json:"35,omitempty"`
	PagesLt  int `json:"36,omitempty"`
	PagesLte int `json:"37,omitempty"`

	Limit    int      `json:"38,omitempty"`
	Offset   int      `json:"39,omitempty"`
	Preloads []string `json:"40,omitempty"`
	Sort     string   `json:"41,omitempty"`
	Order    string   `json:"42,omitempty"`
	All      bool     `json:"43,omitempty"`
}

const (
	orderAsc  = "asc"
	orderDesc = "desc"
)

func (opts *GetArchivesOptions) Validate() {
	opts.Path = strings.ToLower(opts.Path)
	opts.TitleMatch = Slugify(opts.TitleMatch)
	opts.TitleWildcard = Slugify(opts.TitleWildcard)

	opts.ArtistsMatch = SlugifyStrings(opts.ArtistsMatch)
	opts.ArtistsMatchAnd = SlugifyStrings(opts.ArtistsMatchAnd)
	opts.ArtistsWildcard = SlugifyStrings(opts.ArtistsWildcard)
	opts.ArtistsWildcardAnd = SlugifyStrings(opts.ArtistsWildcardAnd)
	opts.ExcludedArtistsMatch = SlugifyStrings(opts.ExcludedArtistsMatch)
	opts.ExcludedArtistsWildcard = SlugifyStrings(opts.ExcludedArtistsWildcard)

	opts.CirclesMatch = SlugifyStrings(opts.CirclesMatch)
	opts.CirclesMatchAnd = SlugifyStrings(opts.CirclesMatchAnd)
	opts.CirclesWildcard = SlugifyStrings(opts.CirclesWildcard)
	opts.CirclesWildcardAnd = SlugifyStrings(opts.CirclesWildcardAnd)
	opts.ExcludedCirclesMatch = SlugifyStrings(opts.ExcludedCirclesMatch)
	opts.ExcludedCirclesWildcard = SlugifyStrings(opts.ExcludedCirclesWildcard)

	opts.MagazinesMatch = SlugifyStrings(opts.MagazinesMatch)
	opts.MagazinesMatchAnd = SlugifyStrings(opts.MagazinesMatchAnd)
	opts.MagazinesWildcard = SlugifyStrings(opts.MagazinesWildcard)
	opts.MagazinesWildcardAnd = SlugifyStrings(opts.MagazinesWildcardAnd)
	opts.ExcludedMagazinesMatch = SlugifyStrings(opts.ExcludedMagazinesMatch)
	opts.ExcludedMagazinesWildcard = SlugifyStrings(opts.ExcludedMagazinesWildcard)

	opts.ParodiesMatch = SlugifyStrings(opts.ParodiesMatch)
	opts.ParodiesMatchAnd = SlugifyStrings(opts.ParodiesMatchAnd)
	opts.ParodiesWildcard = SlugifyStrings(opts.ParodiesWildcard)
	opts.ParodiesWildcardAnd = SlugifyStrings(opts.ParodiesWildcardAnd)
	opts.ExcludedParodiesMatch = SlugifyStrings(opts.ExcludedParodiesMatch)
	opts.ExcludedParodiesWildcard = SlugifyStrings(opts.ExcludedParodiesWildcard)

	opts.TagsMatch = SlugifyStrings(opts.TagsMatch)
	opts.TagsMatchAnd = SlugifyStrings(opts.TagsMatchAnd)
	opts.TagsWildcard = SlugifyStrings(opts.TagsWildcard)
	opts.TagsWildcardAnd = SlugifyStrings(opts.TagsWildcardAnd)
	opts.ExcludedTagsMatch = SlugifyStrings(opts.ExcludedTagsMatch)
	opts.ExcludedTagsWildcard = SlugifyStrings(opts.ExcludedTagsWildcard)

	if !opts.All {
		opts.Limit = Max(opts.Limit, 0)
		opts.Limit = Min(opts.Limit, 100)
		opts.Offset = Max(opts.Offset, 0)
	}

	opts.Preloads = validateArchiveRels(opts.Preloads)
	if strings.EqualFold(opts.Sort, ArchiveCols.ID) {
		opts.Sort = ArchiveCols.ID
	} else if strings.EqualFold(opts.Sort, ArchiveCols.UpdatedAt) {
		opts.Sort = ArchiveCols.UpdatedAt
	} else if strings.EqualFold(opts.Sort, ArchiveCols.PublishedAt) {
		opts.Sort = ArchiveCols.PublishedAt
	} else if strings.EqualFold(opts.Sort, ArchiveCols.Title) {
		opts.Sort = ArchiveCols.Title
	} else if strings.EqualFold(opts.Sort, ArchiveCols.Pages) {
		opts.Sort = ArchiveCols.Pages
	} else {
		opts.Sort = ArchiveCols.CreatedAt
	}

	if strings.EqualFold(opts.Order, orderAsc) {
		opts.Order = orderAsc
	} else {
		opts.Order = orderDesc
	}
}

func (opts *GetArchivesOptions) ToQueries() (selectMods, countMods []QueryMod) {
	countMods = []QueryMod{Select("1")}

	var rawQueries []string
	var rawArgs []any

	if len(opts.Path) > 0 {
		rawQueries = append(rawQueries, "archive.path ILIKE '%' || ? || '%'")
		rawArgs = append(rawArgs, opts.Path)
	}

	if len(opts.TitleMatch) > 0 {
		rawQueries = append(rawQueries, "archive.slug = ?")
		rawArgs = append(rawArgs, opts.TitleMatch)
	} else if len(opts.TitleWildcard) > 0 {
		rawQueries = append(rawQueries, "archive.slug ILIKE '%' || ? || '%'")
		rawArgs = append(rawArgs, opts.TitleWildcard)
	}

	if len(opts.ArtistsMatch) > 0 {
		var q []string
		for _, artist := range opts.ArtistsMatch {
			q = append(q, rawSqlArtistsMatch)
			rawArgs = append(rawArgs, artist)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.ArtistsMatchAnd) > 0 {
		for _, artist := range opts.ArtistsMatchAnd {
			rawQueries = append(rawQueries, rawSqlArtistsMatch)
			rawArgs = append(rawArgs, artist)
		}
	}

	if len(opts.ArtistsWildcard) > 0 {
		var q []string
		for _, artist := range opts.ArtistsWildcard {
			q = append(q, rawSqlArtistsWildcard)
			rawArgs = append(rawArgs, artist)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.ArtistsWildcardAnd) > 0 {
		for _, artist := range opts.ArtistsWildcardAnd {
			rawQueries = append(rawQueries, rawSqlArtistsWildcard)
			rawArgs = append(rawArgs, artist)
		}
	}

	if len(opts.ExcludedArtistsMatch) > 0 {
		for _, artist := range opts.ExcludedArtistsMatch {
			rawQueries = append(rawQueries, rawSqlExcludeArtistsMatch)
			rawArgs = append(rawArgs, artist)
		}
	}

	if len(opts.ExcludedArtistsWildcard) > 0 {
		for _, artist := range opts.ExcludedArtistsWildcard {
			rawQueries = append(rawQueries, rawSqlExcludeArtistsWildcard)
			rawArgs = append(rawArgs, artist)
		}
	}

	if len(opts.CirclesMatch) > 0 {
		var q []string
		for _, circle := range opts.CirclesMatch {
			q = append(q, rawSqlCirclesMatch)
			rawArgs = append(rawArgs, circle)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.CirclesMatchAnd) > 0 {
		for _, circle := range opts.CirclesMatchAnd {
			rawQueries = append(rawQueries, rawSqlCirclesMatch)
			rawArgs = append(rawArgs, circle)
		}
	}

	if len(opts.CirclesWildcard) > 0 {
		var q []string
		for _, circle := range opts.CirclesWildcard {
			q = append(q, rawSqlCirclesWildcard)
			rawArgs = append(rawArgs, circle)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.CirclesWildcardAnd) > 0 {
		for _, circle := range opts.CirclesWildcardAnd {
			rawQueries = append(rawQueries, rawSqlCirclesWildcard)
			rawArgs = append(rawArgs, circle)
		}
	}

	if len(opts.ExcludedCirclesMatch) > 0 {
		for _, circle := range opts.ExcludedCirclesMatch {
			rawQueries = append(rawQueries, rawSqlExcludeCirclesMatch)
			rawArgs = append(rawArgs, circle)
		}
	}

	if len(opts.ExcludedCirclesWildcard) > 0 {
		for _, circle := range opts.ExcludedCirclesWildcard {
			rawQueries = append(rawQueries, rawSqlExcludeCirclesWildcard)
			rawArgs = append(rawArgs, circle)
		}
	}

	if len(opts.MagazinesMatch) > 0 {
		var q []string
		for _, magazine := range opts.MagazinesMatch {
			q = append(q, rawSqlMagazinesMatch)
			rawArgs = append(rawArgs, magazine)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.MagazinesMatchAnd) > 0 {
		for _, magazine := range opts.MagazinesMatchAnd {
			rawQueries = append(rawQueries, rawSqlMagazinesMatch)
			rawArgs = append(rawArgs, magazine)
		}
	}

	if len(opts.MagazinesWildcard) > 0 {
		var q []string
		for _, magazine := range opts.MagazinesWildcard {
			q = append(q, rawSqlMagazinesWildcard)
			rawArgs = append(rawArgs, magazine)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.MagazinesWildcardAnd) > 0 {
		for _, magazine := range opts.MagazinesWildcardAnd {
			rawQueries = append(rawQueries, rawSqlMagazinesWildcard)
			rawArgs = append(rawArgs, magazine)
		}
	}

	if len(opts.ExcludedMagazinesMatch) > 0 {
		for _, magazine := range opts.ExcludedMagazinesMatch {
			rawQueries = append(rawQueries, rawSqlExcludeMagazinesMatch)
			rawArgs = append(rawArgs, magazine)
		}
	}

	if len(opts.ExcludedMagazinesWildcard) > 0 {
		for _, magazine := range opts.ExcludedMagazinesWildcard {
			rawQueries = append(rawQueries, rawSqlExcludeMagazinesWildcard)
			rawArgs = append(rawArgs, magazine)
		}
	}

	if len(opts.ParodiesMatch) > 0 {
		var q []string
		for _, parody := range opts.ParodiesMatch {
			q = append(q, rawSqlParodiesMatch)
			rawArgs = append(rawArgs, parody)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.ParodiesMatchAnd) > 0 {
		for _, parody := range opts.ParodiesMatchAnd {
			rawQueries = append(rawQueries, rawSqlParodiesMatch)
			rawArgs = append(rawArgs, parody)
		}
	}

	if len(opts.ParodiesWildcard) > 0 {
		var q []string
		for _, parody := range opts.ParodiesWildcard {
			q = append(q, rawSqlParodiesWildcard)
			rawArgs = append(rawArgs, parody)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.ParodiesWildcardAnd) > 0 {
		for _, parody := range opts.ParodiesWildcardAnd {
			rawQueries = append(rawQueries, rawSqlParodiesWildcard)
			rawArgs = append(rawArgs, parody)
		}
	}

	if len(opts.ExcludedParodiesMatch) > 0 {
		for _, parody := range opts.ExcludedParodiesMatch {
			rawQueries = append(rawQueries, rawSqlExcludeParodiesMatch)
			rawArgs = append(rawArgs, parody)
		}
	}

	if len(opts.ExcludedParodiesWildcard) > 0 {
		for _, parody := range opts.ExcludedParodiesWildcard {
			rawQueries = append(rawQueries, rawSqlExcludeParodiesWildcard)
			rawArgs = append(rawArgs, parody)
		}
	}

	if len(opts.TagsMatch) > 0 {
		var q []string
		for _, tag := range opts.TagsMatch {
			q = append(q, rawSqlTagsMatch)
			rawArgs = append(rawArgs, tag)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.TagsMatchAnd) > 0 {
		for _, tag := range opts.TagsMatchAnd {
			rawQueries = append(rawQueries, rawSqlTagsMatch)
			rawArgs = append(rawArgs, tag)
		}
	}

	if len(opts.TagsWildcard) > 0 {
		var q []string
		for _, tag := range opts.TagsWildcard {
			q = append(q, rawSqlTagsWildcard)
			rawArgs = append(rawArgs, tag)
		}
		rawQueries = append(rawQueries, JoinOR(q...))
	}

	if len(opts.TagsWildcardAnd) > 0 {
		for _, tag := range opts.TagsWildcardAnd {
			rawQueries = append(rawQueries, rawSqlTagsWildcard)
			rawArgs = append(rawArgs, tag)
		}
	}

	if len(opts.ExcludedTagsMatch) > 0 {
		for _, tag := range opts.ExcludedTagsMatch {
			rawQueries = append(rawQueries, rawSqlExcludeTagsMatch)
			rawArgs = append(rawArgs, tag)
		}
	}

	if len(opts.ExcludedTagsWildcard) > 0 {
		for _, tag := range opts.ExcludedTagsWildcard {
			rawQueries = append(rawQueries, rawSqlExcludeTagsWildcard)
			rawArgs = append(rawArgs, tag)
		}
	}

	if opts.PagesEq > 0 {
		selectMods = append(selectMods, Where("archive.pages = ?", opts.PagesEq))
	} else {
		if opts.PagesGt > 0 {
			selectMods = append(selectMods, Where("archive.pages > ?", opts.PagesGt))
		}
		if opts.PagesGte > 0 {
			selectMods = append(selectMods, Where("archive.pages >= ?", opts.PagesGte))
		}
		if opts.PagesLt > 0 {
			selectMods = append(selectMods, Where("archive.pages < ?", opts.PagesLt))
		}
		if opts.PagesLte > 0 {
			selectMods = append(selectMods, Where("archive.pages <= ?", opts.PagesLte))
		}
	}

	if len(rawQueries) > 0 {
		selectMods = append(selectMods, Where(strings.Join(rawQueries, " AND "), rawArgs...))
	}

	selectMods = append(selectMods, Where("archive.published_at IS NOT NULL AND archive.expunged IS FALSE"))
	countMods = append(countMods, selectMods...)

	selectMods = append(selectMods, OrderBy(fmt.Sprintf("%s %s", opts.Sort, opts.Order)))

	if opts.Limit > 0 {
		selectMods = append(selectMods, Limit(opts.Limit))
	}

	if opts.Offset > 0 {
		selectMods = append(selectMods, Offset(opts.Offset))
	}

	for _, v := range opts.Preloads {
		selectMods = append(selectMods, Load(v))
	}
	return
}

type GetArchivesResult struct {
	Archives []*modext.Archive `json:"data"`
	Total    int               `json:"total"`
	Err      error             `json:"error,omitempty"`
}

func GetArchives(opts *GetArchivesOptions) (result *GetArchivesResult) {
	opts.Validate()

	const prefix = "archives"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Archives.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetArchivesResult)
	}

	result = &GetArchivesResult{Archives: []*modext.Archive{}}
	defer func() {
		if len(result.Archives) > 0 || result.Total > 0 || result.Err != nil {
			cache.Archives.RemoveWithPrefix(prefix, cacheKey)
			cache.Archives.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	selectMods, countMods := opts.ToQueries()
	archives, err := models.Archives(selectMods...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Archives(countMods...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Archives = make([]*modext.Archive, len(archives))
	result.Total = len(count)

	for i, archive := range archives {
		result.Archives[i] = modext.NewArchive(archive).LoadRels(archive)
	}
	return
}

func GetArchiveCount() (int64, error) {
	const archiveCountCacheKey = "archiveCount"
	if c, err := cache.Archives.Get(archiveCountCacheKey); err == nil {
		return c.(int64), nil
	}

	count, err := models.Archives(Where("published_at IS NOT NULL AND expunged IS FALSE")).CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.Unknown
	}

	cache.Archives.Set(archiveCountCacheKey, count, 0)
	return count, nil
}

func GetArchiveStats() (size, pages int64, err error) {
	const (
		archiveSzCacheKey    = "archiveSz"
		archivePagesCacheKey = "archivePages"
	)

	if c, err := cache.Archives.Get(archiveSzCacheKey); err == nil {
		size = c.(int64)
	}

	if c, err := cache.Archives.Get(archivePagesCacheKey); err == nil {
		pages = c.(int64)
	}

	if size > 0 && pages > 0 {
		return
	}

	archives, err := models.Archives(Where("published_at IS NOT NULL AND expunged IS FALSE")).AllG()
	if err != nil {
		log.Println(err)
		err = errs.Unknown
		return
	}

	for _, archive := range archives {
		pages += int64(archive.Pages)
		size += archive.Size
	}

	cache.Archives.Set(archiveSzCacheKey, size, 0)
	cache.Archives.Set(archivePagesCacheKey, pages, 0)
	return
}

func PublishArchive(id int64) (*modext.Archive, error) {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ArchiveNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}

	archive.PublishedAt = null.TimeFrom(time.Now().UTC())
	if err := archive.UpdateG(boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}

	// TODO: Purge cache
	return modext.NewArchive(archive), nil
}

func PublishArchives() error {
	err := models.Archives(Where("published_at IS NULL")).
		UpdateAllG(models.M{"published_at": null.TimeFrom(time.Now().UTC())})
	if err != nil {
		log.Println(err)
		return errs.Unknown
	}
	// TODO: Purge cache
	return nil
}

func UnpublishArchive(id int64) (*modext.Archive, error) {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ArchiveNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}

	archive.PublishedAt.Valid = false
	if err := archive.UpdateG(boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}

	// TODO: Purge cache
	return modext.NewArchive(archive), nil
}

func UnpublishArchives() error {
	err := models.Archives(Where("published_at IS NOT NULL")).
		UpdateAllG(models.M{"published_at": null.NewTime(time.Now(), false)})
	if err != nil {
		log.Println(err)
		return errs.Unknown
	}
	// TODO: Purge cache
	return nil
}

func ExpungeArchive(id int64) error {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ArchiveNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	archive.Expunged = !archive.Expunged
	if err := archive.UpdateG(boil.Whitelist(ArchiveCols.Expunged)); err != nil {
		log.Println(err)
		return errs.Unknown
	}

	return nil
}

func RedirectArchive(from, to int64) error {
	archive, err := models.FindArchiveG(from)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ArchiveNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	archive.RedirectID = null.Int64From(to)
	if err := archive.UpdateG(boil.Whitelist(ArchiveCols.RedirectID)); err != nil {
		log.Println(err)
		return errs.Unknown
	}

	return nil
}

func SetArchiveSource(id int64, source string) error {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ArchiveNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	archive.Source = null.NewString(source, len(source) > 0)
	if err := archive.UpdateG(boil.Whitelist(ArchiveCols.Source)); err != nil {
		log.Println(err)
		return errs.Unknown
	}

	return nil
}

func DeleteArchive(id int64) error {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ArchiveNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	if err := archive.DeleteG(); err != nil {
		log.Println(err)
		return errs.Unknown
	}

	// TODO: Purge cache
	os.Remove(filepath.Join(Config.Directories.Symlinks, strconv.Itoa(int(id))))
	return nil
}

func DeleteArchives() error {
	if err := models.Archives().DeleteAllG(); err != nil {
		log.Println(err)
		return errs.Unknown
	}
	// TODO: Purge cache
	// TODO: Remove symlinks
	return nil
}
