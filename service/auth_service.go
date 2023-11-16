package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
)

type AuthService struct {
	UserRepo         *repo.UserRepo
	LoginAttemptRepo *repo.LoginAttemptRepo
	RefreshTokenRepo *repo.RefreshTokenRepo
}

func NewAuthService(
	userRepo *repo.UserRepo,
	loginAttemptRepo *repo.LoginAttemptRepo,
	refreshTokenRepo *repo.RefreshTokenRepo,
) *AuthService {
	return &AuthService{UserRepo: userRepo, LoginAttemptRepo: loginAttemptRepo, RefreshTokenRepo: refreshTokenRepo}
}

type TokenResponse struct {
	User         model.User `json:"user"`
	AccessToken  string     `json:"accessToken"`
	RefreshToken string     `json:"refreshToken"`
}

func NewTokenResponse(user model.User, accessToken string, refreshToken string) *TokenResponse {
	return &TokenResponse{User: user, AccessToken: accessToken, RefreshToken: refreshToken}
}

func (as AuthService) Login(data form.LoginAttemptData) (*TokenResponse, error) {
	tx, err := as.LoginAttemptRepo.DB.Beginx()
	if err != nil {
		slog.Error("login", "error", err)
		return nil, common.DbInternalError
	}

	userFounded, err := as.UserRepo.GetByEmail(data.Email)
	if err != nil {
		if !errors.Is(repo.NotFoundInRepo, err) {
			slog.Error("login", "error", err)
			return nil, common.DbInternalError
		}
	}
	if userFounded == nil {
		loginAttempt := model.NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		logErr := as.LoginAttemptRepo.Log(loginAttempt, tx)
		if logErr != nil {
			slog.Error("login", "error", "error logging login attempt", "data", data)
			return nil, common.DbInternalError
		}
		slog.Warn("login", "error", "user not found", "data", data)
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("login", "error", "error commiting transaction", "data", data)
			return nil, commitErr
		}
		return nil, ErrorUserNotFound
	}

	if !userFounded.Active {
		loginAttempt := model.NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		logErr := as.LoginAttemptRepo.Log(loginAttempt, tx)
		if logErr != nil {
			slog.Error("login", "error", "error logging login attempt", "data", data)
			return nil, common.DbInternalError
		}
		slog.Warn("login", "error", "user not active", "data", data)
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("login", "error", "error commiting transaction", "data", data)
			return nil, commitErr
		}
		return nil, ErrorUserNotActive
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFounded.Password), []byte(data.Password))
	if err != nil {
		loginAttempt := model.NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		logErr := as.LoginAttemptRepo.Log(loginAttempt, tx)
		if logErr != nil {
			slog.Error("login", "error", "error logging login attempt", "data", data)
			return nil, common.DbInternalError
		}
		slog.Warn("login", "error", "invalid password", "data", data)
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("login", "error", "error commiting transaction", "data", data)
			return nil, commitErr
		}
		return nil, ErrorPasswordIsNotValid
	}

	loginAttempt := model.NewSuccessLoginAttempt(userFounded.Id, data.Email, data.IPAddress, data.UserAgent)
	if err := as.LoginAttemptRepo.Log(loginAttempt, tx); err != nil {
		rollbackErr := tx.Rollback()
		if err != nil {
			slog.Error("login", "error", "error rolling back transaction", "data", data)
			return nil, rollbackErr
		}
		slog.Error("login", "error", "error logging login attempt", "data", data)
		return nil, err
	}

	refreshToken := model.NewRefreshToken(userFounded.Id, uuid.NullUUID{})
	if err := as.RefreshTokenRepo.InvalidateAllRefreshTokens(userFounded.Id, tx); err != nil {
		rollbackErr := tx.Rollback()
		if err != nil {
			slog.Error("login", "error", "error rolling back transaction", "data", data)
			return nil, rollbackErr
		}
		slog.Error("login", "error", "error logging login attempt", "data", data)
		return nil, err
	}
	if err := as.RefreshTokenRepo.SaveRefreshToken(refreshToken, tx); err != nil {
		rollbackErr := tx.Rollback()
		if err != nil {
			slog.Error("login", "error", "error rolling back transaction", "data", data)
			return nil, rollbackErr
		}
		slog.Error("login", "error", "error logging login attempt", "data", data)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("login", "error", "error commiting transaction", "data", data)
		return nil, commitErr
	}

	slog.Info("login", "result", "success", "data", loginAttempt)

	tokenResponse := NewTokenResponse(
		*userFounded,
		GenerateAccessToken(userFounded.Id, userFounded.Role),
		GenerateRefreshToken(userFounded.Id, refreshToken.Id),
	)

	return tokenResponse, nil
}
