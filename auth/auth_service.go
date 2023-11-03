package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/user"
)

var (
	ErrorInvalidRefreshToken = errors.New("invalid refresh token")
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
		err := as.LoginAttemptRepo.Log(loginAttempt, tx)
		if err != nil {
			slog.Error("login", "error", "error logging login attempt", "data", data)
			return nil, common.DbInternalError
		}
		slog.Warn("login", "error", "user not found", "data", data)
		if err := tx.Commit(); err != nil {
			slog.Error("login", "error", "error commiting transaction", "data", data)
			return nil, err
		}
		return nil, user.ErrorUserNotFound
	}

	if !userFounded.Active {
		loginAttempt := NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		err := as.LoginAttemptRepo.Log(loginAttempt, tx)
		if err != nil {
			slog.Error("login", "error", "error logging login attempt", "data", data)
			return nil, common.DbInternalError
		}
		slog.Warn("login", "error", "user not active", "data", data)
		if err := tx.Commit(); err != nil {
			slog.Error("login", "error", "error commiting transaction", "data", data)
			return nil, err
		}
		return nil, user.ErrorUserNotActive
	}

	err := bcrypt.CompareHashAndPassword([]byte(userFounded.Password), []byte(data.Password))
	if err != nil {
		loginAttempt := NewUnSuccessLoginAttempt(data.Email, data.IPAddress, data.UserAgent)
		err := as.LoginAttemptRepo.Log(loginAttempt, tx)
		if err != nil {
			slog.Error("login", "error", "error logging login attempt", "data", data)
			return nil, common.DbInternalError
		}
		slog.Warn("login", "error", "invalid password", "data", data)
		if err := tx.Commit(); err != nil {
			slog.Error("login", "error", "error commiting transaction", "data", data)
			return nil, err
		}
		return nil, err
	}

	loginAttempt := NewSuccessLoginAttempt(userFounded.Id, data.Email, data.IPAddress, data.UserAgent)
	if err := as.LoginAttemptRepo.Log(loginAttempt, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("login", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("login", "error", "error logging login attempt", "data", data)
		return nil, err
	}

	refreshToken := NewRefreshToken(userFounded.Id, uuid.NullUUID{})
	if err := as.RefreshTokenRepo.InvalidateAllRefreshTokens(userFounded.Id, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("login", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("login", "error", "error invalidating refresh tokens", "data", data)
		return nil, err
	}

	if err := as.RefreshTokenRepo.SaveRefreshToken(refreshToken, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("login", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("login", "error", "error saving refresh token", "data", data)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		slog.Error("login", "error", "error commiting transaction", "data", data)
		return nil, err
	}

	tokenResponse := NewTokenResponse(
		*userFounded,
		GenerateAccessToken(userFounded.Id, userFounded.Role),
		GenerateRefreshToken(userFounded.Id, refreshToken.Id),
	)

	slog.Info("login", "result", "success", "data", loginAttempt)

	return tokenResponse, nil
}

func (as AuthService) RotateToken(data *RotateTokenData) (*TokenResponse, error) {
	token, err := ParseToken(data.RefreshToken)
	if err != nil {
		slog.Warn("rotate token", "error", "error parsing token", "data", data)
		return nil, ErrorInvalidRefreshToken
	}

	sub := token.Claims.(jwt.MapClaims)["sub"].(string)
	jti := token.Claims.(jwt.MapClaims)["jti"].(string)
	userId := uuid.MustParse(sub)
	refreshTokenId := uuid.MustParse(jti)

	userFounded, err := as.UserRepo.GetById(userId)
	if err != nil {
		slog.Error("rotate token", "error", "error getting user by id", "userId", userId.String())
		return nil, common.DbInternalError
	}
	if userFounded == nil {
		slog.Warn("rotate token", "error", "user not found", "userId", userId.String())
		return nil, user.ErrorUserNotFound
	}
	if !userFounded.Active {
		slog.Warn("rotate token", "error", "user not active", "userId", userId.String())
		return nil, ErrorInvalidRefreshToken
	}

	refreshTokenFounded, err := as.RefreshTokenRepo.GetById(refreshTokenId)
	if err != nil {
		slog.Error("rotate token", "error", "error getting refresh token by id", "refreshTokenId", refreshTokenId.String())
		return nil, common.DbInternalError
	}

	if !refreshTokenFounded.Active {
		slog.Warn("rotate token", "error", "refresh token not active", "refreshTokenId", refreshTokenId.String())
		return nil, ErrorInvalidRefreshToken
	}

	var tx = as.RefreshTokenRepo.DB.MustBegin()

	refreshToken := NewRefreshTokenWithParent(userFounded.Id, refreshTokenFounded.Id)
	if err := as.RefreshTokenRepo.InvalidateAllRefreshTokens(userFounded.Id, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("rotate token", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("rotate token", "error", "error invalidating refresh tokens", "userId", userId.String())
		return nil, common.DbInternalError
	}

	if err := as.RefreshTokenRepo.SaveRefreshToken(refreshToken, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("rotate token", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("rotate token", "error", "error saving refresh token", "userId", userId.String())
		return nil, common.DbInternalError
	}

	if err := tx.Commit(); err != nil {
		slog.Error("rotate token", "error", "error commiting transaction", "userId", userId.String())
		return nil, err
	}

	tokenResponse := NewTokenResponse(
		*userFounded,
		GenerateAccessToken(userFounded.Id, userFounded.Role),
		GenerateRefreshToken(userFounded.Id, refreshToken.Id),
	)

	slog.Info("rotate token", "result", "success", "userId", sub)

	return tokenResponse, nil

}
