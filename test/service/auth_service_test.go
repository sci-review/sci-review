package test

import (
	"sci-review/auth"
	"sci-review/user"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := user.NewUserRepo(db)
	loginAttemptRepo := auth.NewLoginAttemptRepo(db)
	refreshTokenRepo := auth.NewRefreshTokenRepo(db)
	authService := auth.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"
	newUser := user.NewUser(name, email, password)
	newUser.Active = true
	_ = userRepo.Create(newUser)
	loginAttemptData := auth.LoginAttemptData{
		Email:     email,
		Password:  password,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	tokenResponse, err := authService.Login(loginAttemptData)
	if err != nil {
		t.Error("actual error, expect nil")
	}

	if tokenResponse == nil {
		t.Error("actual nil, expect a valid token response")
	}

	if tokenResponse.AccessToken == "" {
		t.Error("actual empty, expect a valid access token")
	}

	if tokenResponse.RefreshToken == "" {
		t.Error("actual empty, expect a valid refresh token")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := user.NewUserRepo(db)
	loginAttemptRepo := auth.NewLoginAttemptRepo(db)
	refreshTokenRepo := auth.NewRefreshTokenRepo(db)
	authService := auth.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"

	loginAttemptData := auth.LoginAttemptData{
		Email:     email,
		Password:  password,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	tokenResponse, err := authService.Login(loginAttemptData)
	if err == nil {
		t.Error("actual nil, expect error")
	}

	if tokenResponse != nil {
		t.Error("actual token response, expect nil")
	}

	if err != user.ErrorUserNotFound {
		t.Errorf("actual %s, expect %s", err.Error(), user.ErrorUserNotFound.Error())
	}
}

func TestAuthService_Login_UserNotActive(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := user.NewUserRepo(db)
	loginAttemptRepo := auth.NewLoginAttemptRepo(db)
	refreshTokenRepo := auth.NewRefreshTokenRepo(db)
	authService := auth.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"
	newUser := user.NewUser(name, email, password)
	_ = userRepo.Create(newUser)

	loginAttemptData := auth.LoginAttemptData{
		Email:     email,
		Password:  password,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	tokenResponse, err := authService.Login(loginAttemptData)
	if err == nil {
		t.Error("actual nil, expect error")
	}

	if tokenResponse != nil {
		t.Error("actual token response, expect nil")
	}

	if err != user.ErrorUserNotActive {
		t.Errorf("actual %s, expect %s", err.Error(), user.ErrorUserNotActive.Error())
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := user.NewUserRepo(db)
	loginAttemptRepo := auth.NewLoginAttemptRepo(db)
	refreshTokenRepo := auth.NewRefreshTokenRepo(db)
	authService := auth.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"
	newUser := user.NewUser(name, email, password)

	_ = userRepo.Create(newUser)

	wrongPassword := "wrongPassword"
	loginAttemptData := auth.LoginAttemptData{
		Email:     email,
		Password:  wrongPassword,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	tokenResponse, err := authService.Login(loginAttemptData)
	if err == nil {
		t.Error("actual nil, expect error")
	}

	if tokenResponse != nil {
		t.Error("actual token response, expect nil")
	}
}
