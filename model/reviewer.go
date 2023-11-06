package model

import (
	"github.com/google/uuid"
	"time"
)

type Reviewer struct {
	Id           uuid.UUID    `db:"id" json:"id"`
	UserId       uuid.UUID    `db:"user_id" json:"userId"`
	ReviewId     uuid.UUID    `db:"review_id" json:"reviewId"`
	Active       bool         `db:"active" json:"active"`
	ReviewerRole ReviewerRole `db:"role" json:"role"`
	CreatedAt    time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updatedAt"`
}

func NewReviewer(userId uuid.UUID, reviewId uuid.UUID, reviewerRole ReviewerRole) *Reviewer {
	return &Reviewer{
		Id:           uuid.New(),
		UserId:       userId,
		ReviewId:     reviewId,
		Active:       true,
		ReviewerRole: reviewerRole,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
