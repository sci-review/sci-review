package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sci-review/model"
)

type ReviewRepo interface {
	Create(review *model.Review, tx *sqlx.Tx) error
	AddReviewer(reviewer *model.Reviewer, tx *sqlx.Tx) error
	FindAllByUserId(userId uuid.UUID) (*[]model.Review, error)
	FindById(id uuid.UUID) (*model.Review, error)
	GetDB() *sqlx.DB
}

type ReviewRepoSql struct {
	DB *sqlx.DB
}

func NewReviewRepoSql(DB *sqlx.DB) *ReviewRepoSql {
	return &ReviewRepoSql{DB: DB}
}

func (r *ReviewRepoSql) Create(review *model.Review, tx *sqlx.Tx) error {
	query := `
		INSERT INTO reviews (id, owner_id, title, type, start_date, end_date, archived, created_at, updated_at)
		VALUES (:id, :owner_id, :title, :type, :start_date, :end_date, :archived, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, review)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReviewRepoSql) AddReviewer(reviewer *model.Reviewer, tx *sqlx.Tx) error {
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

func (r *ReviewRepoSql) FindAllByUserId(userId uuid.UUID) (*[]model.Review, error) {
	var reviews []model.Review
	query := `
		SELECT r.id, r.owner_id, r.title, r.type, r.start_date, r.end_date, r.archived, r.created_at, r.updated_at
		FROM reviews r
		INNER JOIN reviewers rv ON rv.review_id = r.id
		WHERE rv.user_id = $1
	`
	err := r.DB.Select(&reviews, query, userId)
	if err != nil {
		return nil, err
	}

	if len(reviews) == 0 {
		return &[]model.Review{}, nil
	}
	return &reviews, nil
}

func (r *ReviewRepoSql) FindById(id uuid.UUID) (*model.Review, error) {
	review := model.Review{}
	query := `
		SELECT r.id, r.owner_id, r.title, r.type, r.start_date, r.end_date, r.archived, r.created_at, r.updated_at
		FROM reviews r
		INNER JOIN reviewers rv ON rv.review_id = r.id
		WHERE r.id = $1
	`
	err := r.DB.Get(&review, query, id)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *ReviewRepoSql) GetDB() *sqlx.DB {
	return r.DB
}
