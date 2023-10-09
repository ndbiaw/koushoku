package modext

import "koushoku/models"

type Artist struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count,omitempty" boil:"archive_count"`
}

func NewArtist(model *models.Artist) *Artist {
	if model == nil {
		return nil
	}
	return &Artist{ID: model.ID, Slug: model.Slug, Name: model.Name}
}
