package service

import (
	"errors"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

var (
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorUserNotActive     = errors.New("user not active")
)

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (us *UserService) Create(userCreateForm form.UserCreateForm) (*model.User, error) {
	userFounded, _ := us.UserRepo.GetByEmail(userCreateForm.Email)
	if userFounded != nil {
		slog.Warn("user create", "error", "user already exists", "userData", userCreateForm)
		return nil, ErrorUserAlreadyExists
	}

	user := model.NewUser(userCreateForm.Name, userCreateForm.Email, userCreateForm.Password)
	err := us.UserRepo.Create(user)
	if err != nil {
		slog.Error("user create", "error", err.Error(), "userData", userCreateForm)
		return nil, common.DbInternalError
	}
	slog.Info("user create", "result", "success", "user", user)
	return user, nil
}
