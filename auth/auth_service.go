package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (as AuthService) RotateToken(data *RotateTokenData) (*TokenResponse, error) {
	token, err := ParseToken(data.RefreshToken)
	if err != nil {
		return nil, ErrorInvalidRefreshToken
	}

	sub := token.Claims.(jwt.MapClaims)["sub"].(string)
	jti := token.Claims.(jwt.MapClaims)["jti"].(string)
	userId := uuid.MustParse(sub)
	refreshTokenId := uuid.MustParse(jti)

	userFounded, err := as.UserRepo.GetById(userId)
	if err != nil {
		return nil, common.DbInternalError
	}
	if userFounded == nil {
		return nil, ErrorInvalidRefreshToken
	}
	if !userFounded.Active {
		return nil, ErrorInvalidRefreshToken
	}

	refreshTokenFounded, err := as.RefreshTokenRepo.GetById(refreshTokenId)
	if err != nil {
		return nil, common.DbInternalError
	}

	if !refreshTokenFounded.Active {
		return nil, ErrorInvalidRefreshToken
	}

	var tx = as.RefreshTokenRepo.DB.MustBegin()

	refreshToken := NewRefreshTokenWithParent(userFounded.Id, refreshTokenFounded.Id)
	if err := as.RefreshTokenRepo.InvalidateAllRefreshTokens(userFounded.Id, tx); err != nil {
		tx.Rollback()
		return nil, common.DbInternalError
	}

	if err := as.RefreshTokenRepo.SaveRefreshToken(refreshToken, tx); err != nil {
		tx.Rollback()
		return nil, common.DbInternalError
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
