package model

import (
	"github.com/google/uuid"
	"time"
)

type Investigation struct {
	Id        uuid.UUID           `db:"id" json:"id"`
	UserId    uuid.UUID           `db:"user_id" json:"userId"`
	ReviewId  uuid.UUID           `db:"review_id" json:"reviewId"`
	Question  string              `db:"question" json:"question"`
	Status    InvestigationStatus `db:"status" json:"status"`
	CreatedAt time.Time           `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time           `db:"updated_at" json:"updatedAt"`
}

func NewInvestigation(userId uuid.UUID, reviewId uuid.UUID, question string, status InvestigationStatus) *Investigation {
	return &Investigation{
		Id:        uuid.New(),
		UserId:    userId,
		ReviewId:  reviewId,
		Question:  question,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
