package services

import (
	"log"
	"strings"
	"time"

	"koushoku/cache"
	"koushoku/errs"
	"koushoku/models"
	"koushoku/modext"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	SubmissionCols = models.SubmissionColumns
	SubmissionRels = models.SubmissionRels
)

func CreateSubmission(name, submitter, content string) (*modext.Submission, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, errs.SubmissionNameRequired
	} else if len(name) > 1024 {
		return nil, errs.ArtistNameTooLong
	}

	submitter = strings.TrimSpace(submitter)
	if len(submitter) > 128 {
		return nil, errs.SubmissionSubmitterTooLong
	}

	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil, errs.SubmissionContentRequired
	} else if len(content) > 10240 {
		return nil, errs.SubmissionContentTooLong
	}

	submission := &models.Submission{
		CreatedAt: time.Now().UTC(),
		Name:      name,
		Submitter: null.StringFrom(submitter),
		Content:   content,
	}
	if err := submission.InsertG(boil.Whitelist("created_at", "name", "submitter", "content")); err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewSubmission(submission), nil
}

type GetSubmissionsOptions struct {
	Limit  int `json:"1,omitempty"`
	Offset int `json:"2,omitempty"`
}

type GetSubmissionsResult struct {
	Submissions []*modext.Submission
	Total       int
	Err         error
}

func GetSubmissions(opts GetSubmissionsOptions) (result *GetSubmissionsResult) {
	opts.Limit = Max(opts.Limit, 0)
	opts.Offset = Max(opts.Offset, 0)

	const prefix = "submissions"
	cacheKey := makeCacheKey(opts)
	if c, err := cache.Submissions.GetWithPrefix(prefix, cacheKey); err == nil {
		return c.(*GetSubmissionsResult)
	}

	result = &GetSubmissionsResult{Submissions: []*modext.Submission{}}
	defer func() {
		if len(result.Submissions) > 0 || result.Total > 0 || result.Err != nil {
			cache.Submissions.RemoveWithPrefix(prefix, cacheKey)
			cache.Submissions.SetWithPrefix(prefix, cacheKey, result, 0)
		}
	}()

	selectMods := []QueryMod{Where("accepted = TRUE OR rejected = TRUE")}
	countMods := append([]QueryMod{}, selectMods...)

	if opts.Limit > 0 {
		selectMods = append(selectMods, Limit(opts.Limit))
	}

	if opts.Offset > 0 {
		selectMods = append(selectMods, Offset(opts.Offset))
	}

	selectMods = append(selectMods,
		OrderBy(`COALESCE(accepted_at, rejected_at) DESC, id DESC`),
		Load(SubmissionRels.Archives, OrderBy("id ASC")))
	submissions, err := models.Submissions(selectMods...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := models.Submissions(countMods...).CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Total = int(count)

	result.Submissions = make([]*modext.Submission, len(submissions))
	for i, submission := range submissions {
		result.Submissions[i] = modext.NewSubmission(submission).LoadRels(submission)
	}
	return
}

func AcceptSubmission(id int64, notes string) error {
	submission, err := models.FindSubmissionG(id)
	if err != nil {
		return err
	}

	submission.Accepted = true
	submission.AcceptedAt = null.TimeFrom(time.Now().UTC())
	submission.Notes = null.StringFrom(notes)

	submission.Rejected = false
	submission.RejectedAt.Valid = false

	return submission.UpdateG(boil.Whitelist("accepted", "accepted_at", "rejected", "rejected_at", "notes"))
}

func ListSubmissions() ([]*modext.Submission, error) {
	submissions, err := models.Submissions(OrderBy("id ASC")).AllG()
	if err != nil {
		return nil, err
	}

	result := make([]*modext.Submission, len(submissions))
	for i, submission := range submissions {
		result[i] = modext.NewSubmission(submission)
	}
	return result, nil
}

func RejectSubmission(id int64, note string) error {
	submission, err := models.FindSubmissionG(id)
	if err != nil {
		return err
	}

	submission.Accepted = false
	submission.AcceptedAt.Valid = false

	submission.Rejected = true
	submission.RejectedAt = null.TimeFrom(time.Now().UTC())
	submission.Notes = null.StringFrom(note)

	return submission.UpdateG(boil.Whitelist("accepted", "accepted_at", "rejected", "rejected_at", "notes"))
}

func LinkSubmission(archiveId int64, submissionId int64) error {
	archive, err := models.FindArchiveG(archiveId)
	if err != nil {
		return err
	}

	submission, err := models.FindSubmissionG(submissionId)
	if err != nil {
		return err
	}

	archive.SubmissionID = null.Int64From(submission.ID)
	return archive.UpdateG(boil.Infer())
}
