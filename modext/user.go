package modext

import "koushoku/models"

type User struct {
	ID int64 `json:"id"`

	CreatedAt int64 `json:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt,omitempty"`

	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`

	IsBanned bool `json:"isBanned,omitempty"`
	IsAdmin  bool `json:"isAdmin,omitempty"`

	Favorites []*Archive `json:"favorites,omitempty"`
}

func NewUser(model *models.User) *User {
	if model == nil {
		return nil
	}
	return &User{
		ID: model.ID,

		CreatedAt: model.CreatedAt.Unix(),
		UpdatedAt: model.UpdatedAt.Unix(),

		Email:    model.Email,
		Password: model.Password,
		Name:     model.Name,

		IsBanned: model.IsBanned,
		IsAdmin:  model.IsAdmin,
	}
}

func (user *User) LoadFavorites(model *models.User) *User {
	if model.R == nil || len(model.R.Archives) == 0 {
		return user
	}

	user.Favorites = make([]*Archive, len(model.R.Archives))
	for i, archive := range model.R.Archives {
		user.Favorites[i] = NewArchive(archive)
	}
	return user
}
