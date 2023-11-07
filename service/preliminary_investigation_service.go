package service

import (
	"github.com/google/uuid"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
)

type PreliminaryInvestigationService struct {
	PreliminaryInvestigationRepo *repo.PreliminaryInvestigationRepo
}

func NewPreliminaryInvestigationService(preliminaryInvestigationRepo *repo.PreliminaryInvestigationRepo) *PreliminaryInvestigationService {
	return &PreliminaryInvestigationService{PreliminaryInvestigationRepo: preliminaryInvestigationRepo}
}

func (ps *PreliminaryInvestigationService) Create(data form.PreliminaryInvestigationForm, reviewId uuid.UUID, userId uuid.UUID) (*model.PreliminaryInvestigation, error) {
	preliminaryInvestigation := model.NewPreliminaryInvestigation(userId, reviewId, data.Question)

	err := ps.PreliminaryInvestigationRepo.Create(preliminaryInvestigation)
	if err != nil {
		return nil, err
	}

	return preliminaryInvestigation, nil
}

func (ps *PreliminaryInvestigationService) GetAllByReviewID(reviewId uuid.UUID) ([]model.PreliminaryInvestigation, error) {
	return ps.PreliminaryInvestigationRepo.GetAllByReviewID(reviewId)
}
