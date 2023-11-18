package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	ErrorPasswordIsNotValid    = errors.New("password is not valid")
)

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (us *UserService) Create(userCreateForm form.UserCreateForm) (*model.User, error) {
	userFounded, err := us.UserRepo.GetByEmail(userCreateForm.Email)
	if err != nil {
		if !errors.Is(repo.NotFoundInRepo, err) {
			slog.Error("user create", "error", err)
			return nil, common.DbInternalError
		}
	}
	if userFounded != nil {
		slog.Warn("user create", "error", "user already exists", "userData", userCreateForm)
		return nil, ErrorUserAlreadyExists
	}

	user := model.NewUser(userCreateForm.Name, userCreateForm.Email, userCreateForm.Password)
	err = us.UserRepo.Create(user)
	if err != nil {
		slog.Error("user create", "error", err.Error(), "userData", userCreateForm)
		return nil, common.DbInternalError
	}
	slog.Info("user create", "result", "success", "user", user)
	return user, nil
}

func (us *UserService) FindAll(loggedUserId uuid.UUID) (*[]model.User, error) {
	_, err := us.checkIsAdmin(loggedUserId)
	if err != nil {
		return nil, err
	}

	users, err := us.UserRepo.FindAll()
	if err != nil {
		slog.Warn("user findAll", "error", err.Error())
		return nil, err
	}
	return users, nil
}

func (us *UserService) Activate(loggedUserId uuid.UUID, userId uuid.UUID) error {
	_, err := us.checkIsAdmin(loggedUserId)
	if err != nil {
		return err
	}

	user, err := us.UserRepo.GetById(userId)
	if err != nil {
		if !errors.Is(err, repo.NotFoundInRepo) {
			slog.Warn("user activate", "error", "user not found", "user", userId)
			return common.DbInternalError
		}
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
	_, err := us.checkIsAdmin(loggedUserId)
	if err != nil {
		return err
	}

	user, err := us.UserRepo.GetById(userId)
	if err != nil {
		if !errors.Is(err, repo.NotFoundInRepo) {
			slog.Warn("user deactivate", "error", "user not found", "user", userId)
			return err
		}
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

func (us *UserService) checkIsAdmin(loggedUserId uuid.UUID) (*model.User, error) {
	loggedUser, err := us.UserRepo.GetById(loggedUserId)
	if err != nil {
		if errors.Is(err, repo.NotFoundInRepo) {
			slog.Warn("user check", "error", "logged user not found", "loggedUserId", loggedUserId)
			return nil, ErrorUserNotFound
		} else {
			slog.Warn("user check", "error", err, "loggedUserId", loggedUserId)
			return nil, common.DbInternalError
		}
	}

	if !loggedUser.Active {
		slog.Warn("user check", "error", "logged user is not active", "loggedUserId", loggedUserId)
		return nil, ErrorUserNotActive
	}

	if !loggedUser.IsAdmin() {
		slog.Warn("user check", "error", "logged user is not admin", "loggedUserId", loggedUserId)
		return nil, common.ForbiddenError
	}

	return loggedUser, nil
}

func (us *UserService) CreateAdminUser(name string, email string, password string) error {
	userFounded, err := us.UserRepo.GetByEmail(email)
	if err != nil {
		if !errors.Is(repo.NotFoundInRepo, err) {
			slog.Error("create admin user", "error", err)
			return common.DbInternalError
		}
	}
	if userFounded != nil {
		slog.Info("create admin user", "admin already exists", "email", email)
		return ErrorUserAlreadyExists
	}

	user := model.NewUser(name, email, password)
	user.Role = model.UserAdmin
	user.Active = true
	err = us.UserRepo.Create(user)
	if err != nil {
		slog.Error("create admin user", "error", err.Error(), "email", email)
		return common.DbInternalError
	}
	slog.Info("create admin user", "result", "success", "email", email)
	return nil
}

func (us *UserService) FindById(loggedUserId uuid.UUID, userId uuid.UUID) (*model.User, error) {
	user, err := us.UserRepo.GetById(userId)
	if err != nil {
		if errors.Is(err, repo.NotFoundInRepo) {
			slog.Warn("user find by id", "error", "user not found", "userId", userId)
			return nil, ErrorUserNotFound
		} else {
			slog.Warn("user find by id", "error", err, "userId", userId)
			return nil, common.DbInternalError
		}
	}

	if loggedUserId != userId {
		_, err := us.checkIsAdmin(loggedUserId)
		if err != nil {
			return nil, err
		}
	}

	if !user.Active {
		slog.Warn("change password", "error", "user not active", "userId", userId)
		return nil, ErrorUserNotActive
	}

	return user, nil
}

func (us *UserService) ChangePassword(loggedUserId uuid.UUID, userId uuid.UUID, passwordForm *form.ChangePasswordForm) error {
	user, err := us.UserRepo.GetById(userId)
	if err != nil {
		if errors.Is(err, repo.NotFoundInRepo) {
			slog.Warn("change password", "error", "user not found", "userId", userId)
			return ErrorUserNotFound
		} else {
			slog.Warn("change password", "error", err, "userId", userId)
			return common.DbInternalError
		}
	}

	if loggedUserId != userId {
		return common.ForbiddenError
	}

	if !user.Active {
		slog.Warn("change password", "error", "user not active", "userId", userId)
		return ErrorUserNotActive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordForm.CurrentPassword))
	if err != nil {
		slog.Warn("change password", "error", "invalid password", "userId", userId)
		return ErrorPasswordIsNotValid
	}

	user.SetNewPassword(passwordForm.NewPassword)
	err = us.UserRepo.Update(user)
	if err != nil {
		slog.Warn("change password", "error", err.Error())
		return err
	}

	return nil
}
