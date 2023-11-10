package cache

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
	"sci-review/model"
	"sci-review/repo"
)

type ReviewRepoCache struct {
	ReviewRepo repo.ReviewRepo
	AppCache   *cache.Cache
}

func NewReviewRepoCache(reviewRepo repo.ReviewRepo, appCache *cache.Cache) *ReviewRepoCache {
	return &ReviewRepoCache{ReviewRepo: reviewRepo, AppCache: appCache}
}

func findAllReviewKey(userId uuid.UUID) string {
	return "review:findAll:" + userId.String()
}

func findOneReviewKey(id uuid.UUID) string {
	return "review:findOne:" + id.String()
}

func (r ReviewRepoCache) Create(review *model.Review, tx *sqlx.Tx) error {
	err := r.ReviewRepo.Create(review, tx)
	if err != nil {
		return err
	}

	r.AppCache.Delete(findAllReviewKey(review.OwnerId))
	slog.Debug("ReviewRepoCache.Create: cache cleared", "userId", review.OwnerId)

	return nil
}

func (r ReviewRepoCache) AddReviewer(reviewer *model.Reviewer, tx *sqlx.Tx) error {
	err := r.ReviewRepo.AddReviewer(reviewer, tx)
	if err != nil {
		return err
	}

	r.AppCache.Delete(findAllReviewKey(reviewer.UserId))
	r.AppCache.Delete(findOneReviewKey(reviewer.ReviewId))
	slog.Debug("ReviewRepoCache.AddReviewer: cache cleared", "userId", reviewer.UserId)

	return nil
}

func (r ReviewRepoCache) FindAllByUserId(userId uuid.UUID) (*[]model.Review, error) {
	value, found := r.AppCache.Get(findAllReviewKey(userId))
	if found {
		slog.Debug("ReviewRepoCache.FindAllByUserId: cache hit", "userId", userId)
		return value.(*[]model.Review), nil
	}
	slog.Debug("ReviewRepoCache.FindAllByUserId: cache miss", "userId", userId)

	reviews, err := r.ReviewRepo.FindAllByUserId(userId)
	if err != nil {
		return nil, err
	}

	r.AppCache.Set(findAllReviewKey(userId), reviews, cache.DefaultExpiration)
	slog.Debug("ReviewRepoCache.FindAllByUserId: cache set", "userId", userId)

	return reviews, nil
}

func (r ReviewRepoCache) FindById(id uuid.UUID) (*model.Review, error) {
	value, found := r.AppCache.Get(findOneReviewKey(id))
	if found {
		slog.Debug("ReviewRepoCache.FindById: cache hit", "id", id)
		return value.(*model.Review), nil
	}
	slog.Debug("ReviewRepoCache.FindById: cache miss", "id", id)

	review, err := r.ReviewRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	r.AppCache.Set(findOneReviewKey(id), review, cache.DefaultExpiration)
	slog.Debug("ReviewRepoCache.FindById: cache set", "id", id)

	return review, nil
}

func (r ReviewRepoCache) GetDB() *sqlx.DB {
	return r.ReviewRepo.GetDB()
}
