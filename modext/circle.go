package modext

import "koushoku/models"

type Circle struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count,omitempty" boil:"archive_count"`
}

func NewCircle(model *models.Circle) *Circle {
	if model == nil {
		return nil
	}
	return &Circle{ID: model.ID, Slug: model.Slug, Name: model.Name}
}
