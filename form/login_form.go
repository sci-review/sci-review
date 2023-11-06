package form

import "log/slog"

type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type LoginAttemptData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
}

func (la LoginAttemptData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("email", la.Email),
		slog.String("ipAddress", la.IPAddress),
		slog.String("userAgent", la.UserAgent),
	)
}
