package form

type LogoutForm struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
