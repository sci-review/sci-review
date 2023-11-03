package auth

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepo struct {
	DB *sqlx.DB
}

func NewRefreshTokenRepo(DB *sqlx.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{DB: DB}
}

func (rtr *RefreshTokenRepo) SaveRefreshToken(refreshToken *RefreshToken, tx *sqlx.Tx) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, parent_token_id, issued_at, expires_at, active)
		VALUES (:id, :user_id, :parent_token_id, :issued_at, :expires_at, :active)
	`
	_, err := tx.NamedExec(query, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (rtr *RefreshTokenRepo) InvalidateAllRefreshTokens(userId uuid.UUID, tx *sqlx.Tx) error {
	query := `UPDATE refresh_tokens SET active = false WHERE user_id = $1`

	_, err := tx.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (rtr *RefreshTokenRepo) GetById(id uuid.UUID) (*RefreshToken, error) {
	refreshToken := RefreshToken{}
	err := rtr.DB.Get(&refreshToken, "SELECT * FROM refresh_tokens WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}
