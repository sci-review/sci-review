package auth

import "sci-review/user"

type TokenResponse struct {
	User         user.User `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

func NewTokenResponse(user user.User, accessToken string, refreshToken string) *TokenResponse {
	return &TokenResponse{User: user, AccessToken: accessToken, RefreshToken: refreshToken}
}
