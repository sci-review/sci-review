package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
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

func (pr *PreliminaryInvestigationRepo) GetById(investigationId uuid.UUID) (*model.PreliminaryInvestigation, error) {
	investigation := model.PreliminaryInvestigation{}
	query := `
		SELECT * FROM preliminary_investigations WHERE id = $1
	`
	err := pr.DB.Get(&investigation, query, investigationId)
	if err != nil {
		return nil, err
	}
	return &investigation, nil
}

func (pr *PreliminaryInvestigationRepo) SaveKeyword(investigationKeyword *model.InvestigationKeyword) error {
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

//type Synonyms []string
//
//func (s *Synonyms) Scan(src interface{}) (err error) {
//	var synonyms []string
//	switch src.(type) {
//	case string:
//		// convert this string "{word1,word2,word3}" to []string
//		//remove first and last character
//		src = src.(string)[1 : len(src.(string))-1]
//		// split by comma
//		synonyms = strings.Split(src.(string), ",")
//	case []byte:
//		src = string(src.([]byte))[1 : len(src.([]byte))-1]
//		synonyms = strings.Split(string(src.([]byte)), ",")
//	case nil:
//		synonyms = []string{}
//	default:
//		return errors.New("Incompatible types")
//	}
//	if err != nil {
//		return
//	}
//	*s = synonyms
//	return nil
//}

//type DbKeyword struct {
//	Id              uuid.UUID `db:"id"`
//	UserId          uuid.UUID `db:"user_id"`
//	InvestigationId uuid.UUID `db:"investigation_id"`
//	Word            string    `db:"word"`
//	Synonyms        Synonyms  `db:"synonyms"`
//	CreatedAt       time.Time `db:"created_at"`
//	UpdatedAt       time.Time `db:"updated_at"`
//}

func (pr *PreliminaryInvestigationRepo) GetKeywordsByInvestigationId(investigationId uuid.UUID) ([]model.InvestigationKeyword, error) {
	var keywords []model.InvestigationKeyword
	query := `SELECT * FROM investigation_keywords WHERE investigation_id = $1`
	err := pr.DB.Select(&keywords, query, investigationId)
	if err != nil {
		slog.Error("error", "error", err)
		return nil, err
	}
	return keywords, nil
}
