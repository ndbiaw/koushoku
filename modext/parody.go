package modext

import "koushoku/models"

type Parody struct {
	ID    int64  `json:"id" boil:"id"`
	Slug  string `json:"slug" boil:"slug"`
	Name  string `json:"name" boil:"name"`
	Count int64  `json:"count,omitempty" boil:"archive_count"`
}

func NewParody(model *models.Parody) *Parody {
	if model == nil {
		return nil
	}
	return &Parody{ID: model.ID, Slug: model.Slug, Name: model.Name}
}
