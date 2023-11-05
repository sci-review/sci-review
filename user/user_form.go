package user

import "golang.org/x/exp/slog"

type UserCreateForm struct {
	Name     string `json:"name" form:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" form:"email" validate:"required,email,max=350"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=60"`
}

func (u UserCreateForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", u.Name),
		slog.String("email", u.Email),
	)
}
