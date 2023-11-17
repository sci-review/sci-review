package service

import (
	"github.com/google/uuid"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
	"strings"
)

type InvestigationService struct {
	InvestigationRepo repo.InvestigationRepo
}

func NewInvestigationService(investigationRepo repo.InvestigationRepo) *InvestigationService {
	return &InvestigationService{InvestigationRepo: investigationRepo}
}

func (ps *InvestigationService) Create(data form.InvestigationForm, reviewId uuid.UUID, userId uuid.UUID) (*model.Investigation, error) {
	investigation := model.NewInvestigation(userId, reviewId, data.Question, model.PiStatusInProgress)

	err := ps.InvestigationRepo.Create(investigation)
	if err != nil {
		return nil, err
	}

	return investigation, nil
}

func (ps *InvestigationService) FindAllByReviewID(reviewId uuid.UUID, userId uuid.UUID) (*[]model.Investigation, error) {
	return ps.InvestigationRepo.FindAll(reviewId)
}

func (ps *InvestigationService) FindOneById(investigationId uuid.UUID, userId uuid.UUID) (*model.Investigation, error) {
	return ps.InvestigationRepo.FindOne(investigationId)
}

func (ps *InvestigationService) SaveKeyword(investigationId uuid.UUID, userId uuid.UUID, keywordForm form.KeywordForm) error {
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
	return ps.InvestigationRepo.SaveKeyword(keyword)
}

func (ps *InvestigationService) GetKeywordsByInvestigationId(investigationId uuid.UUID) ([]model.InvestigationKeyword, error) {
	return ps.InvestigationRepo.GetKeywordsByInvestigationId(investigationId)
}
