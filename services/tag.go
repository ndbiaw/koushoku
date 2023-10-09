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

func CreateTag(name string) (*modext.Tag, error) {
	name = strings.Title(strings.TrimSpace(name))
	if len(name) == 0 {
		return nil, errs.TagNameRequired
	} else if len(name) > 128 {
		return nil, errs.TagNameTooLong
	}

	slug := Slugify(name)
	tag, err := models.Tags(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		tag = &models.Tag{Name: name, Slug: slug}
		if err = tag.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.Unknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewTag(tag), nil
}

func GetTag(slug string) (*modext.Tag, error) {
	tag, err := models.Tags(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.TagNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewTag(tag), nil
}

type GetTagsOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetTagsResult struct {
	Tags  []*modext.Tag
	Total int
	Err   error
}

func GetTags(opts GetTagsOptions) (result *GetTagsResult) {
	opts.Limit = Max(opts.Limit, 0)
	opts.Offset = Max(opts.Offset, 0)

	const prefix = "tags"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Taxonomies.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetTagsResult)
	}

	result = &GetTagsResult{Tags: []*modext.Tag{}}
	defer func() {
		if len(result.Tags) > 0 || result.Total > 0 || result.Err != nil {
			cache.Taxonomies.RemoveWithPrefix(prefix, cacheKey)
			cache.Taxonomies.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	q := []QueryMod{
		Select("tag.*", "COUNT(archive.tag_id) AS archive_count"),
		InnerJoin("archive_tags archive ON archive.tag_id = tag.id"),
		GroupBy("tag.id"), OrderBy("tag.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Tags(q...).BindG(context.Background(), &result.Tags)
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Tags().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Total = int(count)
	return
}

func GetTagCount() (int64, error) {
	const cacheKey = "tagCount"
	if c, err := cache.Taxonomies.Get(cacheKey); err == nil {
		return c.(int64), nil
	}

	count, err := models.Tags().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.Unknown
	}

	cache.Taxonomies.Set(cacheKey, count, 0)
	return count, nil
}

var tagIndexes = IndexMap{Cache: make(map[string]bool)}

func IsTagValid(str string) (isValid bool) {
	str = Slugify(str)
	if v, ok := tagIndexes.Get(str); ok {
		return v
	}

	result := GetTags(GetTagsOptions{})
	if result.Err != nil {
		return
	}

	defer tagIndexes.Add(str, isValid)
	for _, tag := range result.Tags {
		if tag.Slug == str {
			isValid = true
			break
		}
	}
	return
}
