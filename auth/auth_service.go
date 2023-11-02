package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"sci-review/user"
)

type AuthService struct {
	UserRepo         *user.UserRepo
	RefreshTokenRepo *RefreshTokenRepo
}

func NewAuthService(userRepo *user.UserRepo, refreshTokenRepo *RefreshTokenRepo) *AuthService {
	return &AuthService{UserRepo: userRepo, RefreshTokenRepo: refreshTokenRepo}
}

func (as AuthService) Login(form LoginForm) (*TokenResponse, error) {
	userFounded, _ := as.UserRepo.GetByEmail(form.Email)
	if userFounded == nil {
		return nil, user.ErrorUserNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(userFounded.Password), []byte(form.Password))
	if err != nil {
		return nil, err
	}

	refreshToken := NewRefreshToken(userFounded.Id, uuid.NullUUID{})

	var tx = as.RefreshTokenRepo.DB.MustBegin()

	if err := as.RefreshTokenRepo.InvalidateAllRefreshTokens(userFounded.Id, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := as.RefreshTokenRepo.SaveRefreshToken(refreshToken, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	tokenResponse := NewTokenResponse(
		*userFounded,
		GenerateAccessToken(userFounded.Id, userFounded.Role),
		GenerateRefreshToken(userFounded.Id, refreshToken.Id),
	)
	return tokenResponse, nil
}
