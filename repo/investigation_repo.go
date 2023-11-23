package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
	"sci-review/model"
)

type InvestigationRepo interface {
	Create(model *model.Investigation) error
	FindAll(reviewID uuid.UUID) (*[]model.Investigation, error)
	FindOne(investigationId uuid.UUID) (*model.Investigation, error)
	SaveKeyword(investigationKeyword *model.InvestigationKeyword) error
	GetKeywordsByInvestigationId(investigationId uuid.UUID) (*[]model.InvestigationKeyword, error)
	GetKeywordsById(id uuid.UUID) (*model.InvestigationKeyword, error)
	DeleteKeyword(id uuid.UUID) error
	UpdateKeyword(investigationKeyword *model.InvestigationKeyword) (*model.InvestigationKeyword, error)
}

type InvestigationRepoSql struct {
	DB *sqlx.DB
}

func NewInvestigationRepoSql(DB *sqlx.DB) *InvestigationRepoSql {
	return &InvestigationRepoSql{DB: DB}
}

func (pr *InvestigationRepoSql) Create(model *model.Investigation) error {
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

func (pr *InvestigationRepoSql) FindAll(reviewID uuid.UUID) (*[]model.Investigation, error) {
	var investigations []model.Investigation
	query := `
		SELECT * FROM investigations WHERE review_id = $1
	`
	err := pr.DB.Select(&investigations, query, reviewID)
	if err != nil {
		return nil, err
	}

	if len(investigations) == 0 {
		return &[]model.Investigation{}, nil
	}

	return &investigations, nil
}

func (pr *InvestigationRepoSql) FindOne(investigationId uuid.UUID) (*model.Investigation, error) {
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

func (pr *InvestigationRepoSql) SaveKeyword(investigationKeyword *model.InvestigationKeyword) error {
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

func (pr *InvestigationRepoSql) GetKeywordsByInvestigationId(investigationId uuid.UUID) (*[]model.InvestigationKeyword, error) {
	var keywords []model.InvestigationKeyword
	query := `SELECT * FROM investigation_keywords WHERE investigation_id = $1`
	err := pr.DB.Select(&keywords, query, investigationId)
	if err != nil {
		slog.Error("error", "error", err)
		return nil, err
	}

	if len(keywords) == 0 {
		return &[]model.InvestigationKeyword{}, nil
	}

	return &keywords, nil
}

func (pr *InvestigationRepoSql) GetKeywordsById(id uuid.UUID) (*model.InvestigationKeyword, error) {
	var keyword model.InvestigationKeyword
	query := `SELECT * FROM investigation_keywords WHERE id = $1`
	err := pr.DB.Get(&keyword, query, id)
	if err != nil {
		return nil, err
	}
	return &keyword, nil

}

func (pr *InvestigationRepoSql) DeleteKeyword(id uuid.UUID) error {
	query := `DELETE FROM investigation_keywords WHERE id = $1`
	_, err := pr.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (pr *InvestigationRepoSql) UpdateKeyword(investigationKeyword *model.InvestigationKeyword) (*model.InvestigationKeyword, error) {
	query := `
		UPDATE investigation_keywords SET word = :word, synonyms = :synonyms, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := pr.DB.NamedExec(query, investigationKeyword)
	if err != nil {
		return nil, err
	}
	return investigationKeyword, nil
}
