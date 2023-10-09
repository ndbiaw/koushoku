package modext

import "koushoku/models"

type Submission struct {
	ID int64 `json:"id"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`

	Name      string `json:"-"`
	Submitter string `json:"submitter,omitempty"`
	Content   string `json:"-"`
	Notes     string `json:"-"`

	AcceptedAt int64 `json:"acceptedAt,omitempty"`
	RejectedAt int64 `json:"rejectedAt,omitempty"`

	Accepted bool `json:"accepted,omitempty"`
	Rejected bool `json:"rejected,omitempty"`

	Archives []*Archive `json:"archives,omitempty"`
}

func NewSubmission(model *models.Submission) *Submission {
	if model == nil {
		return nil
	}

	submission := &Submission{
		ID:        model.ID,
		CreatedAt: model.CreatedAt.Unix(),
		UpdatedAt: model.UpdatedAt.Unix(),

		Name:      model.Name,
		Submitter: model.Submitter.String,
		Content:   model.Content,
		Notes:     model.Notes.String,

		Accepted: model.Accepted,
		Rejected: model.Rejected,
	}

	if model.AcceptedAt.Valid {
		submission.AcceptedAt = model.AcceptedAt.Time.Unix()
	}

	if model.RejectedAt.Valid {
		submission.RejectedAt = model.RejectedAt.Time.Unix()
	}

	return submission
}

func (submission *Submission) LoadRels(model *models.Submission) *Submission {
	if model == nil || model.R == nil {
		return submission
	}

	submission.LoadArchives(model)

	return submission
}

func (submission *Submission) LoadArchives(model *models.Submission) *Submission {
	if model == nil || model.R == nil || len(model.R.Archives) == 0 {
		return submission
	}

	submission.Archives = make([]*Archive, len(model.R.Archives))
	for i, archive := range model.R.Archives {
		submission.Archives[i] = NewArchive(archive)
	}

	return submission
}
