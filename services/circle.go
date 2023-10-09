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

func CreateCircle(name string) (*modext.Circle, error) {
	name = strings.Title(strings.TrimSpace(name))
	if len(name) == 0 {
		return nil, errs.CircleNameRequired
	} else if len(name) > 128 {
		return nil, errs.CircleNameTooLong
	}

	slug := Slugify(name)
	circle, err := models.Circles(Where("slug = ?", slug)).OneG()
	if err == sql.ErrNoRows {
		circle = &models.Circle{Name: name, Slug: slug}
		if err = circle.InsertG(boil.Infer()); err != nil {
			log.Println(err)
			return nil, errs.Unknown
		}
	} else if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewCircle(circle), nil
}

func GetCircle(slug string) (*modext.Circle, error) {
	circle, err := models.Circles(Where("slug = ?", slug)).OneG()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.CircleNotFound
		}
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewCircle(circle), nil
}

type GetCirclesOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetCirclesResult struct {
	Circles []*modext.Circle
	Total   int
	Err     error
}

func GetCircles(opts GetCirclesOptions) (result *GetCirclesResult) {
	opts.Limit = Max(opts.Limit, 0)
	opts.Offset = Max(opts.Offset, 0)

	const prefix = "circles"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Taxonomies.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetCirclesResult)
	}

	result = &GetCirclesResult{Circles: []*modext.Circle{}}
	defer func() {
		if len(result.Circles) > 0 || result.Total > 0 || result.Err != nil {
			cache.Taxonomies.RemoveWithPrefix(prefix, cacheKey)
			cache.Taxonomies.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	q := []QueryMod{
		Select("circle.*", "COUNT(archive.circle_id) AS archive_count"),
		InnerJoin("archive_circles archive ON archive.circle_id = circle.id"),
		GroupBy("circle.id"), OrderBy("circle.name ASC"),
	}

	if opts.Limit > 0 {
		q = append(q, Limit(opts.Limit))
		if opts.Offset > 0 {
			q = append(q, Offset(opts.Offset))
		}
	}

	err := models.Circles(q...).BindG(context.Background(), &result.Circles)
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Circles().CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Total = int(count)
	return
}

func GetCircleCount() (int64, error) {
	const cacheKey = "circleCount"
	if c, err := cache.Taxonomies.Get(cacheKey); err == nil {
		return c.(int64), nil
	}

	count, err := models.Circles().CountG()
	if err != nil {
		log.Println(err)
		return 0, errs.Unknown
	}

	cache.Taxonomies.Set(cacheKey, count, 0)
	return count, nil
}

var circleIndexes = IndexMap{Cache: make(map[string]bool)}

func IsCircleValid(str string) (isValid bool) {
	str = Slugify(str)
	if v, ok := circleIndexes.Get(str); ok {
		return v
	}

	result := GetCircles(GetCirclesOptions{})
	if result.Err != nil {
		return
	}

	defer circleIndexes.Add(str, isValid)
	for _, circle := range result.Circles {
		if circle.Slug == str {
			isValid = true
			break
		}
	}
	return
}
