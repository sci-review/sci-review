package test

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
)

func ClearTables() {
	db := GetDb()
	db.MustExec("DELETE FROM login_attempts")
	db.MustExec("DELETE FROM refresh_tokens")
	db.MustExec("DELETE FROM users")
}

func GetDb() *sqlx.DB {
	dataSourceName := "postgresql://postgres:postgres@localhost:5432/sci_review"
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return db
}
