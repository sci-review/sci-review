package user

import "errors"

type UserService struct {
	UserRepo *UserRepo
}

var ErrorUserAlreadyExists = errors.New("user already exists")

func NewUserService(userRepo *UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (us *UserService) Create(userCreateForm UserCreateForm) (*User, error) {
	userFounded, _ := us.UserRepo.GetByEmail(userCreateForm.Email)
	if userFounded != nil {
		return nil, ErrorUserAlreadyExists
	}

	user := NewUser(userCreateForm.Name, userCreateForm.Email, userCreateForm.Password)
	err := us.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
