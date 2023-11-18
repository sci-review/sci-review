package form

import "golang.org/x/exp/slog"

type ChangePasswordForm struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=6,max=60"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=60"`
}

func (u ChangePasswordForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("currentPassword", "******"),
		slog.String("newPassword", "******"),
	)
}
