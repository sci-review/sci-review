package model

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	Id         uuid.UUID  `db:"id" json:"id"`
	Title      string     `db:"title" json:"title"`
	ReviewType ReviewType `db:"type" json:"type"`
	StartDate  time.Time  `db:"start_date" json:"startDate"`
	EndDate    time.Time  `db:"end_date" json:"endDate"`
	Archived   bool       `db:"archived" json:"archived"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
	Reviewers  []Reviewer `db:"-" json:"reviewers"`
}

func NewReview(title string, reviewType ReviewType, startDate time.Time, endDate time.Time) *Review {
	return &Review{
		Id:         uuid.New(),
		Title:      title,
		ReviewType: reviewType,
		StartDate:  startDate,
		EndDate:    endDate,
		Archived:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
