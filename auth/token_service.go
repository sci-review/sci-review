package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"os"
	"sci-review/user"
	"strconv"
	"time"
)

func GenerateAccessToken(userId uuid.UUID, role user.Role) string {
	now := time.Now()
	key := os.Getenv("JWT_KEY")
	iss := os.Getenv("JWT_ISSUER")
	accessTokenDuration, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))
	iat := now.Unix()
	exp := now.Add(time.Duration(accessTokenDuration) * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  iss,
			"sub":  userId.String(),
			"iat":  iat,
			"exp":  exp,
			"role": role,
		})
	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		slog.Warn("generate access token", "error", "error signing token", "userId", userId)
		return ""
	}
	slog.Info("generate access token", "result", "success", "userId", userId)

	return signedString
}

func GenerateRefreshToken(userId uuid.UUID, refreshTokenId uuid.UUID) string {
	now := time.Now()
	key := os.Getenv("JWT_KEY")
	iss := os.Getenv("JWT_ISSUER")
	refreshTokenDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))
	iat := now.Unix()
	exp := now.Add(time.Duration(refreshTokenDuration) * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": iss,
			"sub": userId.String(),
			"iat": iat,
			"exp": exp,
			"jti": refreshTokenId.String(),
		})
	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		slog.Warn("generate refresh token", "error", "error signing token", "userId", userId)
		return ""
	}

	slog.Info("generate refresh token", "result", "success", "userId", userId)
	return signedString
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	key := os.Getenv("JWT_KEY")
	iss := os.Getenv("JWT_ISSUER")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("parse token", "error", "unexpected signing method", "token", tokenString)
			return nil, jwt.ErrSignatureInvalid
		}

		if token.Claims.(jwt.MapClaims)["iss"] != iss {
			slog.Error("parse token", "error", "invalid issuer", "token", tokenString)
			return nil, jwt.ErrTokenUnverifiable
		}

		return []byte(key), nil
	})

	if err != nil {
		slog.Error("parse token", "error", err.Error())
		return nil, err
	}

	return token, nil
}
