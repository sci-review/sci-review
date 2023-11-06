package repo

import (
	"github.com/jmoiron/sqlx"
	"sci-review/model"
)

type ReviewRepo struct {
	DB *sqlx.DB
}

func NewReviewRepo(DB *sqlx.DB) *ReviewRepo {
	return &ReviewRepo{DB: DB}
}

func (r *ReviewRepo) Create(review *model.Review, tx *sqlx.Tx) error {
	query := `
		INSERT INTO reviews (id, title, type, start_date, end_date, archived, created_at, updated_at)
		VALUES (:id, :title, :type, :start_date, :end_date, :archived, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, review)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepo) AddReviewer(reviewer *model.Reviewer, tx *sqlx.Tx) error {
	query := `
		INSERT INTO reviewers (id, user_id, review_id, role, active, created_at, updated_at)
		VALUES (:id, :user_id, :review_id, :role, :active, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, reviewer)
	if err != nil {
		return err
	}
	return nil
}
