package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sci-review/model"
)

type PreliminaryInvestigationRepo struct {
	DB *sqlx.DB
}

func NewPreliminaryInvestigationRepo(DB *sqlx.DB) *PreliminaryInvestigationRepo {
	return &PreliminaryInvestigationRepo{DB: DB}
}

func (pr *PreliminaryInvestigationRepo) Create(model *model.PreliminaryInvestigation) error {
	query := `
		INSERT INTO preliminary_investigations (id, user_id, review_id, question, status, created_at, updated_at)
		VALUES (:id, :user_id, :review_id, :question, :status, :created_at, :updated_at)
	`
	_, err := pr.DB.NamedExec(query, model)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PreliminaryInvestigationRepo) GetAllByReviewID(reviewID uuid.UUID) ([]model.PreliminaryInvestigation, error) {
	var models []model.PreliminaryInvestigation
	query := `
		SELECT * FROM preliminary_investigations WHERE review_id = $1
	`
	err := pr.DB.Select(&models, query, reviewID)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (pr *PreliminaryInvestigationRepo) GetById(investigationId uuid.UUID) (model.PreliminaryInvestigation, error) {
	investigation := model.PreliminaryInvestigation{}
	query := `
		SELECT * FROM preliminary_investigations WHERE id = $1
	`
	err := pr.DB.Get(&investigation, query, investigationId)
	if err != nil {
		return investigation, err
	}
	return investigation, nil
}
