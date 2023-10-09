package main

import (
	"fmt"
	"net/http"
	"strings"

	"koushoku/server"
	"koushoku/services"
)

const (
	archiveTmplName = "archive.html"
	readerTmplName  = "reader.html"
)

func archive(c *server.Context) {
	if c.TryCache(archiveTmplName) {
		return
	}

	id, err := c.ParamInt64("id")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	result := services.GetArchive(id, services.GetArchiveOptions{
		Preloads: []string{
			services.ArchiveRels.Artists,
			services.ArchiveRels.Circles,
			services.ArchiveRels.Magazines,
			services.ArchiveRels.Parodies,
			services.ArchiveRels.Tags,
			services.ArchiveRels.Submission,
		},
	})
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	if (result.Archive.RedirectId > 0) && (result.Archive.RedirectId != id) {
		c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d", result.Archive.RedirectId))
		return
	}

	slug := c.Param("slug")
	isJson := strings.HasSuffix(slug, ".json")
	if isJson {
		slug = strings.TrimSuffix(slug, ".json")
	}

	if !strings.EqualFold(slug, result.Archive.Slug) {
		slug = result.Archive.Slug
		if isJson {
			slug += ".json"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s", result.Archive.ID, slug))
		return
	}

	if isJson {
		c.JSON(http.StatusOK, result.Archive)
	} else {
		c.SetData("archive", result.Archive)
		c.Cache(http.StatusOK, archiveTmplName)
	}
}

func read(c *server.Context) {
	if c.TryCache(readerTmplName) {
		return
	}

	id, err := c.ParamInt64("id")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	pageNum, err := c.ParamInt("pageNum")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	result := services.GetArchive(id, services.GetArchiveOptions{})
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	slug := c.Param("slug")
	if !strings.EqualFold(slug, result.Archive.Slug) {
		if pageNum <= 0 || int16(pageNum) > result.Archive.Pages {
			c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s/1", result.Archive.ID, result.Archive.Slug))
		} else {
			c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s/%d", result.Archive.ID, result.Archive.Slug, pageNum))
		}
		return
	}

	if pageNum <= 0 || int16(pageNum) > result.Archive.Pages {
		c.Redirect(http.StatusFound, fmt.Sprintf("/archive/%d/%s/1", id, result.Archive.Slug))
		return
	}

	c.SetData("archive", result.Archive)
	c.SetData("pageNum", pageNum)
	c.Cache(http.StatusOK, readerTmplName)
}
