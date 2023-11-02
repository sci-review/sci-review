package auth

import (
	"github.com/jmoiron/sqlx"
)

type LoginAttemptRepo struct {
	DB *sqlx.DB
}

func NewLoginAttemptRepo(DB *sqlx.DB) *LoginAttemptRepo {
	return &LoginAttemptRepo{DB: DB}
}

func (lar LoginAttemptRepo) Log(loginAttempt *LoginAttempt, tx *sqlx.Tx) error {
	query := `
		INSERT INTO login_attempts (id, user_id, email, success, ip_address, user_agent, timestamp)
		VALUES (:id, :user_id, :email, :success, :ip_address, :user_agent, :timestamp)
	`
	if tx != nil {
		_, err := tx.NamedExec(query, loginAttempt)
		if err != nil {
			return err
		}
	} else {
		_, err := lar.DB.NamedExec(query, loginAttempt)
		if err != nil {
			return err
		}
	}
	return nil
}
