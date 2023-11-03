package user

import (
	"errors"
	"golang.org/x/exp/slog"
	"sci-review/common"
)

type UserService struct {
	UserRepo *UserRepo
}

var (
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorUserNotActive     = errors.New("user not active")
)

func NewUserService(userRepo *UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (us *UserService) Create(userCreateForm UserCreateForm) (*User, error) {
	userFounded, _ := us.UserRepo.GetByEmail(userCreateForm.Email)
	if userFounded != nil {
		slog.Warn("user create", "error", "user already exists", "userData", userCreateForm)
		return nil, ErrorUserAlreadyExists
	}

	user := NewUser(userCreateForm.Name, userCreateForm.Email, userCreateForm.Password)
	err := us.UserRepo.Create(user)
	if err != nil {
		slog.Error("user create", "error", err.Error(), "userData", userCreateForm)
		return nil, common.DbInternalError
	}
	slog.Info("user create", "result", "success", "user", user)
	return user, nil
}
