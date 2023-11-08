package service

import (
	"github.com/google/uuid"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
	"strings"
)

type PreliminaryInvestigationService struct {
	PreliminaryInvestigationRepo *repo.PreliminaryInvestigationRepo
}

func NewPreliminaryInvestigationService(preliminaryInvestigationRepo *repo.PreliminaryInvestigationRepo) *PreliminaryInvestigationService {
	return &PreliminaryInvestigationService{PreliminaryInvestigationRepo: preliminaryInvestigationRepo}
}

func (ps *PreliminaryInvestigationService) Create(data form.PreliminaryInvestigationForm, reviewId uuid.UUID, userId uuid.UUID) (*model.PreliminaryInvestigation, error) {
	preliminaryInvestigation := model.NewPreliminaryInvestigation(userId, reviewId, data.Question, model.PiStatusInProgress)

	err := ps.PreliminaryInvestigationRepo.Create(preliminaryInvestigation)
	if err != nil {
		return nil, err
	}

	return preliminaryInvestigation, nil
}

func (ps *PreliminaryInvestigationService) GetAllByReviewID(reviewId uuid.UUID) ([]model.PreliminaryInvestigation, error) {
	return ps.PreliminaryInvestigationRepo.GetAllByReviewID(reviewId)
}

func (ps *PreliminaryInvestigationService) GetById(investigationId uuid.UUID, userId uuid.UUID) (model.PreliminaryInvestigation, error) {
	return ps.PreliminaryInvestigationRepo.GetById(investigationId)
}

func (ps *PreliminaryInvestigationService) SaveKeyword(investigationId uuid.UUID, userId uuid.UUID, keywordForm form.KeywordForm) error {
	formSynonyms := strings.Split(keywordForm.Synonyms, "\n")
	var synonyms []string

	// trim spaces and remove empty synonyms
	for _, synonym := range formSynonyms {
		synonym = strings.TrimSpace(synonym)
		if synonym != "" {
			synonyms = append(synonyms, synonym)
		}
	}

	keyword := model.NewInvestigationKeyword(userId, investigationId, keywordForm.Word, synonyms)
	return ps.PreliminaryInvestigationRepo.SaveKeyword(keyword)
}

func (ps *PreliminaryInvestigationService) GetKeywordsByInvestigationId(investigationId uuid.UUID) ([]model.InvestigationKeyword, error) {
	return ps.PreliminaryInvestigationRepo.GetKeywordsByInvestigationId(investigationId)
}
