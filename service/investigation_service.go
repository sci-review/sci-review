package service

import (
	"github.com/google/uuid"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
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

func (ps *InvestigationService) SaveKeyword(investigationId uuid.UUID, userId uuid.UUID, keywordForm form.KeywordForm) (*model.InvestigationKeyword, error) {
	keyword := model.NewInvestigationKeyword(userId, investigationId, keywordForm.Word, keywordForm.Synonyms)
	err := ps.InvestigationRepo.SaveKeyword(keyword)
	if err != nil {
		return nil, err
	}

	return keyword, nil
}

func (ps *InvestigationService) GetKeywordsByInvestigationId(investigationId uuid.UUID) (*[]model.InvestigationKeyword, error) {
	return ps.InvestigationRepo.GetKeywordsByInvestigationId(investigationId)
}

func (ps *InvestigationService) DeleteKeyword(keywordId uuid.UUID) error {
	return ps.InvestigationRepo.DeleteKeyword(keywordId)
}

func (ps *InvestigationService) UpdateKeyword(keywordId uuid.UUID, keywordForm form.KeywordForm) (*model.InvestigationKeyword, error) {
	keyword, err := ps.InvestigationRepo.GetKeywordsById(keywordId)
	if err != nil {
		return nil, err
	}

	keyword.Update(keywordForm.Word, keywordForm.Synonyms)
	updatedKeyword, err := ps.InvestigationRepo.UpdateKeyword(keyword)
	if err != nil {
		return nil, err
	}

	return updatedKeyword, nil
}
