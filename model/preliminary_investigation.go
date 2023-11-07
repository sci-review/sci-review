package model

import (
	"github.com/google/uuid"
	"time"
)

type PreliminaryInvestigation struct {
	Id        uuid.UUID `db:"id" json:"id"`
	UserId    uuid.UUID `db:"user_id" json:"userId"`
	ReviewId  uuid.UUID `db:"review_id" json:"reviewId"`
	Question  string    `db:"question" json:"question"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func NewPreliminaryInvestigation(userId uuid.UUID, reviewId uuid.UUID, question string) *PreliminaryInvestigation {
	return &PreliminaryInvestigation{
		Id:        uuid.New(),
		UserId:    userId,
		ReviewId:  reviewId,
		Question:  question,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
