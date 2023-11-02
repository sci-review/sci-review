package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"sci-review/user"
)

type AuthService struct {
	UserRepo         *user.UserRepo
	LoginAttemptRepo *LoginAttemptRepo
	RefreshTokenRepo *RefreshTokenRepo
}

func NewAuthService(userRepo *user.UserRepo, loginAttemptRepo *LoginAttemptRepo, refreshTokenRepo *RefreshTokenRepo) *AuthService {
	return &AuthService{UserRepo: userRepo, LoginAttemptRepo: loginAttemptRepo, RefreshTokenRepo: refreshTokenRepo}
}

func (as AuthService) Login(data LoginAttemptData) (*TokenResponse, error) {
	var tx = as.RefreshTokenRepo.DB.MustBegin()

	userFounded, _ := as.UserRepo.GetByEmail(data.Email)
	if userFounded == nil {
		loginAttempt := NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		as.LoginAttemptRepo.Log(loginAttempt, tx)
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, user.ErrorUserNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(userFounded.Password), []byte(data.Password))
	if err != nil {
		loginAttempt := NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		as.LoginAttemptRepo.Log(loginAttempt, tx)
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, err
	}

	loginAttempt := NewSuccessLoginAttempt(userFounded.Id, data.Email, data.IPAddress, data.UserAgent)
	if err := as.LoginAttemptRepo.Log(loginAttempt, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	refreshToken := NewRefreshToken(userFounded.Id, uuid.NullUUID{})
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
