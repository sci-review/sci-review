package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sci-review/model"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(DB *sqlx.DB) *UserRepo {
	return &UserRepo{DB: DB}
}

func (ur *UserRepo) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, password, role, active, created_at, updated_at)
		VALUES (:id, :name, :email, :password, :role, :active, :created_at, :updated_at)
	`
	_, err := ur.DB.NamedExec(query, user)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetByEmail(email string) (*model.User, error) {
	user := model.User{}
	err := ur.DB.Get(&user, "SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundInRepo
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) GetById(id uuid.UUID) (*model.User, error) {
	user := model.User{}
	err := ur.DB.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundInRepo
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) FindAll() (*[]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users ORDER by created_at DESC`
	err := ur.DB.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (ur *UserRepo) Update(user *model.User) error {
	query := `
		UPDATE users 
		SET 
			name = :name, 
			email = :email, 
			password = :password,
			role = :role, 
			active = :active, 
			updated_at = :updated_at 
		WHERE id = :id;
	`
	_, err := ur.DB.NamedExec(query, user)
	if err != nil {
		return err
	}
	return nil
}
