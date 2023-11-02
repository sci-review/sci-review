package auth

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Id            uuid.UUID     `db:"id" json:"id"`
	UserId        uuid.UUID     `db:"user_id" json:"userId"`
	ParentTokenId uuid.NullUUID `db:"parent_token_id" json:"parentTokenId"`
	IssuedAt      time.Time     `db:"issued_at" json:"issuedAt"`
	ExpiresAt     time.Time     `db:"expires_at" json:"expiresAt"`
	Active        bool          `db:"active" json:"active"`
}

func NewRefreshToken(userId uuid.UUID, parentTokenId uuid.NullUUID) *RefreshToken {
	return &RefreshToken{
		Id:            uuid.New(),
		UserId:        userId,
		ParentTokenId: parentTokenId,
		IssuedAt:      time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour * 24),
		Active:        true,
	}
}
