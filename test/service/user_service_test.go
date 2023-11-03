package test

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"sci-review/user"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	ClearTables()
	userRepo := user.NewUserRepo(GetDb())
	userService := user.NewUserService(userRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"

	userCreateForm := user.UserCreateForm{
		Name:     name,
		Email:    email,
		Password: password,
	}

	userCreated, err := userService.Create(userCreateForm)

	if err != nil {
		t.Error(err.Error())
		return
	}
	if userCreated.Name != userCreateForm.Name {
		t.Error("user name not match")
		return
	}
	if userCreated.Email != userCreateForm.Email {
		t.Error("user email not match")
		return
	}
}

func TestUserService_Create_UserAlreadyExists(t *testing.T) {
	ClearTables()
	userRepo := user.NewUserRepo(GetDb())
	userService := user.NewUserService(userRepo)

	name := "Test test"
	email := "teste@email.com"
	password := "test123"

	_ = userRepo.Create(user.NewUser(name, email, password))

	userCreateForm := user.UserCreateForm{
		Name:     name,
		Email:    email,
		Password: password,
	}

	userCreated, err := userService.Create(userCreateForm)

	if userCreated != nil {
		t.Error("actual user, expect nil")
	}

	if err == nil {
		t.Error("actual nil, expect error")
	}

	if err != user.ErrorUserAlreadyExists {
		t.Errorf("actual %s, expect %s", err.Error(), user.ErrorUserAlreadyExists.Error())
	}
}
