package modext

import "koushoku/models"

type Magazine struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count,omitempty" boil:"archive_count"`
}

func NewMagazine(model *models.Magazine) *Magazine {
	if model == nil {
		return nil
	}
	return &Magazine{ID: model.ID, Slug: model.Slug, Name: model.Name}
}
