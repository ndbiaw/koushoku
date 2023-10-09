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

func CreateParody(name string) (*modext.Parody, error) {
	name = strings.Title(strings.TrimSpace(name))
	if len(name) == 0 {
		return nil, errs.ParodyNameRequired
	} else if len(name) > 128 {
		return nil, errs.ParodyNameTooLong
	}

	slug := Slugify(name)
	parody, err := models.Parodies(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		parody = &models.Parody{Name: name, Slug: slug}
		if err = parody.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.Unknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewParody(parody), nil
}

func GetParody(slug string) (*modext.Parody, error) {
	parody, err := models.Parodies(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ParodyNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewParody(parody), nil
}

type GetParodiesOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetParodiesResult struct {
	Parodies []*modext.Parody
	Total    int
	Err      error
}

func GetParodies(opts GetParodiesOptions) (result *GetParodiesResult) {
	opts.Limit = Max(opts.Limit, 0)
	opts.Offset = Max(opts.Offset, 0)

	const prefix = "parodies"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Taxonomies.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetParodiesResult)
	}

	result = &GetParodiesResult{Parodies: []*modext.Parody{}}
	defer func() {
		if len(result.Parodies) > 0 || result.Total > 0 || result.Err != nil {
			cache.Taxonomies.RemoveWithPrefix(prefix, cacheKey)
			cache.Taxonomies.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	q := []QueryMod{
		Select("parody.*", "COUNT(archive.parody_id) AS archive_count"),
		InnerJoin("archive_parodies archive ON archive.parody_id = parody.id"),
		GroupBy("parody.id"), OrderBy("parody.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Parodies(q...).BindG(context.Background(), &result.Parodies)
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Parodies().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Total = int(count)
	return
}

func GetParodyCount() (int64, error) {
	const cacheKey = "parodyCount"
	if c, err := cache.Taxonomies.Get(cacheKey); err == nil {
		return c.(int64), nil
	}

	count, err := models.Parodies().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.Unknown
	}

	cache.Taxonomies.Set(cacheKey, count, 0)
	return count, nil
}

var parodyIndexes = IndexMap{Cache: make(map[string]bool)}

func IsParodyValid(str string) (isValid bool) {
	str = Slugify(str)
	if v, ok := parodyIndexes.Get(str); ok {
		return v
	}

	result := GetParodies(GetParodiesOptions{})
	if result.Err != nil {
		return
	}

	defer parodyIndexes.Add(str, isValid)
	for _, parody := range result.Parodies {
		if parody.Slug == str {
			isValid = true
			break
		}
	}
	return
}
