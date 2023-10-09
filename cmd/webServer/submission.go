package main

import (
	"fmt"
	"koushoku/server"
	"koushoku/services"
	"math"
	"net/http"
	"strconv"
)

const (
	submitTmplName      = "submit.html"
	submissionsTmplname = "submissions.html"
)

func submit(c *server.Context) {
	if !c.TryCache(submitTmplName) {
		c.Cache(http.StatusOK, submitTmplName)
	}
}

type SubmitPayload struct {
	Name      string `form:"name"`
	Submitter string `form:"submitter"`
	Content   string `form:"content"`
}

func submitPost(c *server.Context) {
	payload := &SubmitPayload{}
	c.Bind(payload)

	_, err := services.CreateSubmission(payload.Name, payload.Submitter, payload.Content)
	if err != nil {
		c.SetData("lastSubmissionName", payload.Name)
		c.SetData("lastSubmissionSubmitter", payload.Submitter)
		c.SetData("lastSubmissionContent", payload.Content)
		c.SetData("error", err)
		c.HTML(http.StatusBadRequest, submitTmplName)
		return
	}
	c.SetData("message", "Your submission has been submitted.")
	c.HTML(http.StatusOK, submitTmplName)
}

func submisisions(c *server.Context) {
	if c.TryCache(submissionsTmplname) {
		return
	}

	q := &SearchQueries{}
	c.BindQuery(q)

	page, _ := strconv.Atoi(c.Query("page"))
	result := services.GetSubmissions(services.GetSubmissionsOptions{
		Limit:  listingLimit,
		Offset: listingLimit * (page - 1),
	})
	if result.Err != nil {
		c.SetData("error", result.Err)
		c.HTML(http.StatusInternalServerError, "error.html")
		return
	}

	c.SetData("page", page)
	if page > 0 {
		c.SetData("name", fmt.Sprintf("Submissions: Page %d", page))
	} else {
		c.SetData("name", "Submissions")
	}

	totalPages := int(math.Ceil(float64(result.Total) / float64(listingLimit)))
	c.SetData("data", result.Submissions)
	c.SetData("total", result.Total)
	c.SetData("pagination", services.CreatePagination(page, totalPages))

	c.Cache(http.StatusOK, submissionsTmplname)
}
