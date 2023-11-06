package test

import (
	"sci-review/form"
	"sci-review/middleware"
	"sci-review/model"
	"sci-review/repo"
	"sci-review/service"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := repo.NewUserRepo(db)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	refreshTokenRepo := middleware.NewRefreshTokenRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"
	newUser := model.NewUser(name, email, password)
	newUser.Active = true
	_ = userRepo.Create(newUser)
	loginAttemptData := form.LoginAttemptData{
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
	userRepo := repo.NewUserRepo(db)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	refreshTokenRepo := middleware.NewRefreshTokenRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"

	loginAttemptData := form.LoginAttemptData{
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

	if err != service.ErrorUserNotFound {
		t.Errorf("actual %s, expect %s", err.Error(), service.ErrorUserNotFound.Error())
	}
}

func TestAuthService_Login_UserNotActive(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := repo.NewUserRepo(db)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	refreshTokenRepo := middleware.NewRefreshTokenRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"
	newUser := model.NewUser(name, email, password)
	_ = userRepo.Create(newUser)

	loginAttemptData := form.LoginAttemptData{
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

	if err != service.ErrorUserNotActive {
		t.Errorf("actual %s, expect %s", err.Error(), service.ErrorUserNotActive.Error())
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	ClearTables()
	db := GetDb()
	userRepo := repo.NewUserRepo(db)
	loginAttemptRepo := repo.NewLoginAttemptRepo(db)
	refreshTokenRepo := middleware.NewRefreshTokenRepo(db)
	authService := service.NewAuthService(userRepo, loginAttemptRepo, refreshTokenRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"
	ipAddress := "127.0.0.1"
	userAgent := "chrome"
	newUser := model.NewUser(name, email, password)

	_ = userRepo.Create(newUser)

	wrongPassword := "wrongPassword"
	loginAttemptData := form.LoginAttemptData{
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
