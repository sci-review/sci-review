package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Role      Role      `db:"role" json:"role"`
	Active    bool      `db:"active" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewUser(name string, email string, password string) *User {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{
		Id:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  string(passwordHash),
		Role:      Reviewer,
		Active:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
