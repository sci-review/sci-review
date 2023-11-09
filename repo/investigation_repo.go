package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
	"sci-review/model"
)

type InvestigationRepo struct {
	DB *sqlx.DB
}

func NewInvestigationRepo(DB *sqlx.DB) *InvestigationRepo {
	return &InvestigationRepo{DB: DB}
}

func (pr *InvestigationRepo) Create(model *model.Investigation) error {
	query := `
		INSERT INTO investigations (id, user_id, review_id, question, status, created_at, updated_at)
		VALUES (:id, :user_id, :review_id, :question, :status, :created_at, :updated_at)
	`
	_, err := pr.DB.NamedExec(query, model)
	if err != nil {
		return err
	}
	return nil
}

func (pr *InvestigationRepo) GetAllByReviewID(reviewID uuid.UUID) ([]model.Investigation, error) {
	var models []model.Investigation
	query := `
		SELECT * FROM investigations WHERE review_id = $1
	`
	err := pr.DB.Select(&models, query, reviewID)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (pr *InvestigationRepo) GetById(investigationId uuid.UUID) (*model.Investigation, error) {
	investigation := model.Investigation{}
	query := `
		SELECT * FROM investigations WHERE id = $1
	`
	err := pr.DB.Get(&investigation, query, investigationId)
	if err != nil {
		return nil, err
	}
	return &investigation, nil
}

func (pr *InvestigationRepo) SaveKeyword(investigationKeyword *model.InvestigationKeyword) error {
	query := `
		INSERT INTO investigation_keywords (id, user_id, investigation_id, word, synonyms, created_at, updated_at)
		VALUES (:id, :user_id, :investigation_id, :word, :synonyms, :created_at, :updated_at)
	`
	_, err := pr.DB.NamedExec(query, investigationKeyword)
	if err != nil {
		return err
	}
	return nil
}

func (pr *InvestigationRepo) GetKeywordsByInvestigationId(investigationId uuid.UUID) ([]model.InvestigationKeyword, error) {
	var keywords []model.InvestigationKeyword
	query := `SELECT * FROM investigation_keywords WHERE investigation_id = $1`
	err := pr.DB.Select(&keywords, query, investigationId)
	if err != nil {
		slog.Error("error", "error", err)
		return nil, err
	}
	return keywords, nil
}
