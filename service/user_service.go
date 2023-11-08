package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
	"time"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

var (
	ErrorUserAlreadyExists     = errors.New("user already exists")
	ErrorUserNotFound          = errors.New("user not found")
	ErrorUserNotActive         = errors.New("user not active")
	ErrorUserAlreadyActive     = errors.New("user already active")
	ErrorUserAlreadyDeactivate = errors.New("user already deactivate")
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

func (us *UserService) FindAll(loggedUserId uuid.UUID) (*[]model.User, error) {
	loggedUser, err := us.UserRepo.GetById(loggedUserId)
	if err != nil {
		slog.Warn("user findall", "error", "logged user not found", "loggedUserId", loggedUserId)
		return nil, err
	}

	if !loggedUser.IsAdmin() {
		slog.Warn("user findall", "error", "logged user is not admin", "loggedUserId", loggedUserId)
		return nil, common.ForbiddenError
	}

	users, err := us.UserRepo.FindAll()
	if err != nil {
		slog.Warn("user findAll", "error", err.Error())
		return nil, err
	}
	return users, nil
}

func (us *UserService) Activate(loggedUserId uuid.UUID, userId uuid.UUID) error {
	loggedUser, err := us.UserRepo.GetById(loggedUserId)
	if err != nil {
		slog.Warn("user activate", "error", "logged user not found", "loggedUserId", loggedUserId)
		return err
	}

	if !loggedUser.IsAdmin() {
		slog.Warn("user activate", "error", "logged user is not admin", "loggedUserId", loggedUserId)
		return common.ForbiddenError
	}

	user, err := us.UserRepo.GetById(userId)
	if err != nil {
		slog.Warn("user activate", "error", "user not found", "user", userId)
		return err
	}

	if user.Active {
		slog.Warn("user activate", "error", "user already active", "user", userId)
		return ErrorUserAlreadyActive
	}

	user.Active = true
	user.UpdatedAt = time.Now()
	err = us.UserRepo.Update(user)
	if err != nil {
		slog.Warn("user activate", "error", err.Error())
		return err
	}

	return nil
}

func (us *UserService) Deactivate(loggedUserId uuid.UUID, userId uuid.UUID) error {
	loggedUser, err := us.UserRepo.GetById(loggedUserId)
	if err != nil {
		slog.Warn("user deactivate", "error", "logged user not found", "loggedUserId", loggedUserId)
		return err
	}

	if !loggedUser.IsAdmin() {
		slog.Warn("user deactivate", "error", "logged user is not admin", "loggedUserId", loggedUserId)
		return common.ForbiddenError
	}

	user, err := us.UserRepo.GetById(userId)
	if err != nil {
		slog.Warn("user deactivate", "error", "user not found", "user", userId)
		return err
	}

	if !user.Active {
		slog.Warn("user deactivate", "error", "user already deactive", "user", userId)
		return ErrorUserAlreadyDeactivate
	}

	user.Active = false
	user.UpdatedAt = time.Now()
	err = us.UserRepo.Update(user)
	if err != nil {
		slog.Warn("user deactivate", "error", err.Error())
		return err
	}

	return nil
}
