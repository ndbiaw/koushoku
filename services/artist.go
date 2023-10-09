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

func CreateArtist(name string) (*modext.Artist, error) {
	name = strings.Title(strings.TrimSpace(name))
	if len(name) == 0 {
		return nil, errs.ArtistNameRequired
	} else if len(name) > 128 {
		return nil, errs.ArtistNameTooLong
	}

	slug := Slugify(name)
	artist, err := models.Artists(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		artist = &models.Artist{Name: name, Slug: slug}
		if err = artist.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.Unknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}

	return modext.NewArtist(artist), nil
}

func GetArtist(slug string) (*modext.Artist, error) {
	artist, err := models.Artists(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ArtistNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewArtist(artist), nil
}

type GetArtistsOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetArtistsResult struct {
	Artists []*modext.Artist
	Total   int
	Err     error
}

func GetArtists(opts GetArtistsOptions) (result *GetArtistsResult) {
	opts.Limit = Max(opts.Limit, 0)
	opts.Offset = Max(opts.Offset, 0)

	const prefix = "artists"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Taxonomies.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetArtistsResult)
	}

	result = &GetArtistsResult{Artists: []*modext.Artist{}}
	defer func() {
		if len(result.Artists) > 0 || result.Total > 0 || result.Err != nil {
			cache.Taxonomies.RemoveWithPrefix(prefix, cacheKey)
			cache.Taxonomies.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	q := []QueryMod{
		Select("artist.*", "COUNT(archive.artist_id) AS archive_count"),
		InnerJoin("archive_artists archive ON archive.artist_id = artist.id"),
		GroupBy("artist.id"), OrderBy("artist.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Artists(q...).BindG(context.Background(), &result.Artists)
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Artists().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Total = int(count)
	return
}

func GetArtistCount() (int64, error) {
	const cachekey = "artistCount"
	if c, err := cache.Taxonomies.Get(cachekey); err == nil {
		return c.(int64), nil
	}

	count, err := models.Artists().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.Unknown
	}

	cache.Taxonomies.Set(cachekey, count, 0)
	return count, nil
}

var artistIndexes = IndexMap{Cache: make(map[string]bool)}

func IsArtistValid(str string) (isValid bool) {
	str = Slugify(str)
	if v, ok := artistIndexes.Get(str); ok {
		return v
	}

	result := GetArtists(GetArtistsOptions{})
	if result.Err != nil {
		return
	}

	defer artistIndexes.Add(str, isValid)
	for _, artist := range result.Artists {
		if artist.Slug == str {
			isValid = true
			break
		}
	}
	return
}
