package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		fmt.Println(err)
		return ""
	}

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
		fmt.Println(err)
		return ""
	}

	return signedString
}
