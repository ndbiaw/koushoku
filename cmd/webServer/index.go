package main

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"koushoku/server"
	"koushoku/services"
)

type SearchQueries struct {
	Search string `form:"q"`
	Page   int    `form:"page"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

const (
	indexLimit      = 25
	indexTmplName   = "index.html"
	aboutTmplName   = "about.html"
	statsTmplName   = "stats.html"
	searchTmplName  = "search.html"
	sitemapTmplName = "sitemap.xml"
)

func createNewSearchQueries(c *server.Context) *SearchQueries {
	q := &SearchQueries{}
	c.BindQuery(q)

	q.Search = strings.TrimSpace(q.Search)
	if len(q.Sort) == 0 {
		q.Sort = "created_at"
	} else {
		q.Sort = strings.ToLower(q.Sort)
	}

	if len(q.Order) == 0 {
		q.Order = "desc"
	} else {
		q.Order = strings.ToLower(q.Order)
	}
	return q
}

func index(c *server.Context) {
	if c.TryCache(indexTmplName) {
		return
	}

	q := createNewSearchQueries(c)
	result := services.GetArchives(&services.GetArchivesOptions{
		Limit:  indexLimit,
		Offset: indexLimit * (q.Page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circles,
			services.ArchiveRels.Magazines,
			services.ArchiveRels.Tags,
		},
		Sort:  q.Sort,
		Order: q.Order,
	})
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("queries", q)
	if q.Page > 0 {
		c.SetData("name", fmt.Sprintf("Browse: Page %d", q.Page))
	} else {
		c.SetData("name", "Home")
	}

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)
	c.SetData("pagination", services.CreatePagination(q.Page, totalPages))

	c.Cache(http.StatusOK, indexTmplName)
}

var rgx = regexp.MustCompile(`(?i)-?(artist|circle|magazine|parody|tag|title|pages)(&|\|)?(\*)?:(<=?|>=?)?(\".*?\"|[^\s]+)`)

const (
	opAnd      = "&"
	opOr       = "|"
	opWildcard = "*"
)

func search(c *server.Context) {
	if c.TryCache(searchTmplName) {
		return
	}

	q := createNewSearchQueries(c)
	opts := &services.GetArchivesOptions{
		Limit:  indexLimit,
		Offset: indexLimit * (q.Page - 1),
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circles,
			services.ArchiveRels.Magazines,
			services.ArchiveRels.Tags,
		},
		Sort:  q.Sort,
		Order: q.Order,
	}

	if len(q.Search) > 0 {
		matches := rgx.FindAllStringSubmatch(q.Search, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				isExclude := strings.HasPrefix(match[0], "-")
				isAnd := strings.EqualFold(strings.TrimSpace(match[2]), opAnd)
				isWildcard := strings.EqualFold(strings.TrimSpace(match[3]), opWildcard)

				value := match[5]
				values := strings.Split(value, ",")
				numOp := strings.TrimSpace(match[4])

				switch match[1] {
				case "artist":
					for _, v := range values {
						v = strings.TrimSpace(v)
						if isExclude {
							if isWildcard {
								opts.ExcludedArtistsWildcard = append(opts.ExcludedArtistsWildcard, v)
							} else {
								opts.ExcludedArtistsMatch = append(opts.ExcludedArtistsMatch, v)
							}
						} else if isAnd {
							if isWildcard {
								opts.ArtistsWildcardAnd = append(opts.ArtistsWildcardAnd, v)
							} else {
								opts.ArtistsMatchAnd = append(opts.ArtistsMatchAnd, v)
							}
						} else if isWildcard {
							opts.ArtistsWildcard = append(opts.ArtistsWildcard, v)
						} else {
							opts.ArtistsMatch = append(opts.ArtistsMatch, v)
						}
					}
				case "circle":
					for _, v := range values {
						v = strings.TrimSpace(v)
						if isExclude {
							if isWildcard {
								opts.ExcludedCirclesWildcard = append(opts.ExcludedCirclesWildcard, v)
							} else {
								opts.ExcludedCirclesMatch = append(opts.ExcludedCirclesMatch, v)
							}
						} else if isAnd {
							if isWildcard {
								opts.CirclesWildcardAnd = append(opts.CirclesWildcardAnd, v)
							} else {
								opts.CirclesMatchAnd = append(opts.CirclesMatchAnd, v)
							}
						} else if isWildcard {
							opts.CirclesWildcard = append(opts.CirclesWildcard, v)
						} else {
							opts.CirclesMatch = append(opts.CirclesMatch, v)
						}
					}
				case "magazine":
					for _, v := range values {
						v = strings.TrimSpace(v)
						if isExclude {
							if isWildcard {
								opts.ExcludedMagazinesWildcard = append(opts.ExcludedMagazinesWildcard, v)
							} else {
								opts.ExcludedMagazinesMatch = append(opts.ExcludedMagazinesMatch, v)
							}
						} else if isAnd {
							if isWildcard {
								opts.MagazinesWildcardAnd = append(opts.MagazinesWildcardAnd, v)
							} else {
								opts.MagazinesMatchAnd = append(opts.MagazinesMatchAnd, v)
							}
						} else if isWildcard {
							opts.MagazinesWildcard = append(opts.MagazinesWildcard, v)
						} else {
							opts.MagazinesMatch = append(opts.MagazinesMatch, v)
						}
					}
				case "pages":
					n, _ := strconv.Atoi(value)
					if strings.EqualFold(numOp, ">") {
						opts.PagesGt = n
					} else if strings.EqualFold(numOp, ">=") {
						opts.PagesGte = n
					} else if strings.EqualFold(numOp, "<") {
						opts.PagesLt = n
					} else if strings.EqualFold(numOp, "<=") {
						opts.PagesLte = n
					} else {
						opts.PagesEq = n
					}
				case "parody":
					for _, v := range values {
						v = strings.TrimSpace(v)
						if isExclude {
							if isWildcard {
								opts.ExcludedParodiesWildcard = append(opts.ExcludedParodiesWildcard, v)
							} else {
								opts.ExcludedParodiesMatch = append(opts.ExcludedParodiesMatch, v)
							}
						} else if isAnd {
							if isWildcard {
								opts.ParodiesWildcardAnd = append(opts.ParodiesWildcardAnd, v)
							} else {
								opts.ParodiesMatchAnd = append(opts.ParodiesMatchAnd, v)
							}
						} else if isWildcard {
							opts.ParodiesWildcard = append(opts.ParodiesWildcard, v)
						} else {
							opts.ParodiesMatch = append(opts.ParodiesMatch, v)
						}
					}
				case "tag":
					for _, v := range values {
						v = strings.TrimSpace(v)
						if isExclude {
							if isWildcard {
								opts.ExcludedTagsWildcard = append(opts.ExcludedTagsWildcard, v)
							} else {
								opts.ExcludedTagsMatch = append(opts.ExcludedTagsMatch, v)
							}
						} else if isAnd {
							if isWildcard {
								opts.TagsWildcardAnd = append(opts.TagsWildcardAnd, v)
							} else {
								opts.TagsMatchAnd = append(opts.TagsMatchAnd, v)
							}
						} else if isWildcard {
							opts.TagsWildcard = append(opts.TagsWildcard, v)
						} else {
							opts.TagsMatch = append(opts.TagsMatch, v)
						}
					}
				case "title":
					if isWildcard {
						opts.TitleWildcard = strings.Join(values, ",")
					} else {
						opts.TitleMatch = strings.Join(values, ",")
					}
				}
			}
		} else {
			if services.IsArtistValid(q.Search) {
				opts.ArtistsMatch = append(opts.ArtistsMatch, q.Search)
			}

			if services.IsCircleValid(q.Search) {
				opts.CirclesMatch = append(opts.CirclesMatch, q.Search)
			}

			if services.IsParodyValid(q.Search) {
				opts.ParodiesMatch = append(opts.ParodiesMatch, q.Search)
			}

			if services.IsTagValid(q.Search) {
				opts.TagsMatch = append(opts.TagsMatch, q.Search)
			} else {
				arr := strings.Split(q.Search, " ")
				if len(arr) > 1 {
					for _, v := range arr {
						if services.IsTagValid(v) {
							opts.TagsMatch = append(opts.TagsMatch, v)
						}
					}
				}
			}

			if len(opts.ArtistsMatch) == 0 && len(opts.CirclesMatch) == 0 &&
				len(opts.MagazinesMatch) == 0 && len(opts.ParodiesMatch) == 0 &&
				len(opts.TagsMatch) == 0 {
				opts.Path = q.Search
			}
		}
	}

	result := services.GetArchives(opts)
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("queries", q)
	hasQueries := len(q.Search) > 0
	c.SetData("hasQueries", hasQueries)

	if hasQueries {
		c.SetData("name", fmt.Sprintf("Search: %s", q.Search))
	} else {
		c.SetData("name", "Browse")
	}

	c.SetData("archives", result.Archives)
	c.SetData("total", result.Total)

	totalPages := int(math.Ceil(float64(result.Total) / float64(indexLimit)))
	c.SetData("pagination", services.CreatePagination(q.Page, totalPages))

	if len(result.Archives) > 0 {
		c.Cache(http.StatusOK, searchTmplName)
	} else {
		c.Cache(http.StatusNotFound, searchTmplName)
	}
}

func about(c *server.Context) {
	if !c.TryCache(aboutTmplName) {
		c.Cache(http.StatusOK, aboutTmplName)
	}
}

func stats(c *server.Context) {
	if !c.TryCache(statsTmplName) {
		c.SetData("stats", services.GetStats())
		c.Cache(http.StatusOK, statsTmplName)
	}
}

func sitemap(c *server.Context) {
	if c.TryCache(sitemapTmplName) {
		return
	}

	archives := services.GetArchives(&services.GetArchivesOptions{Order: "published_at", All: true})
	if archives.Err != nil {
		c.ErrorJSON(http.StatusInternalServerError, "Failed to get archives", archives.Err)
		return
	}

	artists := services.GetArtists(services.GetArtistsOptions{Limit: 10000})
	if artists.Err != nil {
		c.ErrorJSON(http.StatusInternalServerError, "Failed to get artists", artists.Err)
		return
	}

	magazines := services.GetMagazines(services.GetMagazinesOptions{Limit: 10000})
	if magazines.Err != nil {
		c.ErrorJSON(http.StatusInternalServerError, "Failed to get magazines", magazines.Err)
		return
	}

	tags := services.GetTags(services.GetTagsOptions{Limit: 10000})
	if tags.Err != nil {
		c.ErrorJSON(http.StatusInternalServerError, "Failed to get tags", tags.Err)
		return
	}

	c.SetData("archives", archives.Archives)
	c.SetData("artists", artists.Artists)
	c.SetData("magazines", magazines.Magazines)
	c.SetData("tags", tags.Tags)
	c.Cache(http.StatusOK, sitemapTmplName)
}
