package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
	"time"
)

type ReviewService struct {
	ReviewRepo repo.ReviewRepo
	UserRepo   *repo.UserRepo
}

func NewReviewService(reviewRepo repo.ReviewRepo, userRepo *repo.UserRepo) *ReviewService {
	return &ReviewService{ReviewRepo: reviewRepo, UserRepo: userRepo}
}

var (
	ErrorParseStartDate = errors.New("start date must be in format YYYY-MM-DD")
	ErrorParseEndDate   = errors.New("end date must be in format YYYY-MM-DD")
	ErrorReviewDate     = errors.New("end date must be after start date")
	ErrorReviewNotFound = errors.New("review not found")
)

func (s *ReviewService) Create(data form.ReviewCreateForm, userId uuid.UUID) (*model.Review, error) {
	startDate, err := time.Parse(time.RFC3339, data.StartDate)
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrorParseStartDate
	}

	endDate, err := time.Parse(time.RFC3339, data.EndDate)
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrorParseEndDate
	}

	if endDate.Before(startDate) || endDate.Equal(startDate) {
		slog.Error("end date must be after start date")
		return nil, ErrorReviewDate
	}

	user, err := s.UserRepo.GetById(userId)
	if err != nil {
		return nil, err
	}

	if !user.Active {
		slog.Error("create review", "error", "user not active", "user", user)
		return nil, ErrorUserNotActive
	}

	tx := s.ReviewRepo.GetDB().MustBegin()
	defer tx.Rollback()

	review := model.NewReview(userId, data.Title, data.Type, startDate, endDate)
	err = s.ReviewRepo.Create(review, tx)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	reviewer := model.NewReviewer(userId, review.Id, model.ReviewerOwner)
	err = s.ReviewRepo.AddReviewer(reviewer, tx)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) FindAll(userId uuid.UUID) (*[]model.Review, error) {
	return s.ReviewRepo.FindAllByUserId(userId)
}

func (s *ReviewService) FindById(id uuid.UUID, userId uuid.UUID) (*model.Review, error) {
	review, err := s.ReviewRepo.FindById(id)
	if err != nil {
		return nil, ErrorReviewNotFound
	}
	return review, nil
}
