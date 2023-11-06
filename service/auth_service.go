package service

import (
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
}

func NewAuthService(userRepo *repo.UserRepo, loginAttemptRepo *repo.LoginAttemptRepo) *AuthService {
	return &AuthService{UserRepo: userRepo, LoginAttemptRepo: loginAttemptRepo}
}

func (as AuthService) Login(data form.LoginAttemptData) (*model.User, error) {
	var tx = as.LoginAttemptRepo.DB.MustBegin()

	userFounded, _ := as.UserRepo.GetByEmail(data.Email)
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

	err := bcrypt.CompareHashAndPassword([]byte(userFounded.Password), []byte(data.Password))
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
		return nil, err
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

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("login", "error", "error commiting transaction", "data", data)
		return nil, commitErr
	}

	slog.Info("login", "result", "success", "data", loginAttempt)

	return userFounded, nil
}
