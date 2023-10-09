package services

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"koushoku/cache"
	"koushoku/errs"
	"koushoku/models"
	"koushoku/modext"

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CreateMagazine(name string) (*modext.Magazine, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, errs.MagazineNameRequired
	} else if len(name) > 128 {
		return nil, errs.MagazineNameTooLong
	}

	slug := Slugify(name)
	magazine, err := models.Magazines(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		magazine = &models.Magazine{Name: name, Slug: slug}
		if err = magazine.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.Unknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewMagazine(magazine), nil
}

func GetMagazine(slug string) (*modext.Magazine, error) {
	magazine, err := models.Magazines(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.MagazineNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewMagazine(magazine), nil
}

type GetMagazinesOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetMagazinesResult struct {
	Magazines []*modext.Magazine
	Total     int
	Err       error
}

func GetMagazines(opts GetMagazinesOptions) (result *GetMagazinesResult) {
	opts.Limit = Max(opts.Limit, 0)
	opts.Offset = Max(opts.Offset, 0)

	const prefix = "magazines"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Taxonomies.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetMagazinesResult)
	}

	result = &GetMagazinesResult{Magazines: []*modext.Magazine{}}
	defer func() {
		if len(result.Magazines) > 0 || result.Total > 0 || result.Err != nil {
			cache.Taxonomies.RemoveWithPrefix(prefix, cacheKey)
			cache.Taxonomies.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	q := []QueryMod{
		Select("magazine.*", "COUNT(archive.magazine_id) AS archive_count"),
		InnerJoin("archive_magazines archive ON archive.magazine_id = magazine.id"),
		GroupBy("magazine.id"), OrderBy("magazine.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Magazines(q...).BindG(context.Background(), &result.Magazines)
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Magazines().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
	}

	result.Total = int(count)
	return
}

func GetMagazineCount() (int64, error) {
	const cachekey = "magazineCount"
	if c, err := cache.Taxonomies.Get(cachekey); err == nil {
		return c.(int64), nil
	}

	count, err := models.Magazines().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.Unknown
	}

	cache.Taxonomies.Set(cachekey, count, 0)
	return count, nil
}

var magazineIndexes = IndexMap{Cache: make(map[string]bool)}

func IsMagazineValid(str string) (isValid bool) {
	str = Slugify(str)
	if v, ok := magazineIndexes.Get(str); ok {
		return v
	}

	result := GetMagazines(GetMagazinesOptions{})
	if result.Err != nil {
		return
	}

	defer magazineIndexes.Add(str, isValid)
	for _, magazine := range result.Magazines {
		if magazine.Slug == str {
			isValid = true
			break
		}
	}
	return
}
