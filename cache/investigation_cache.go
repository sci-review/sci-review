package cache

import (
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
	"sci-review/model"
	"sci-review/repo"
)

func findAllInvestigationKey(review uuid.UUID) string {
	return "investigation:findAll:" + review.String()
}

func findOneInvestigationKey(id uuid.UUID) string {
	return "investigation:findOne:" + id.String()
}

type InvestigationRepoCache struct {
	InvestigationRepo repo.InvestigationRepo
	AppCache          *cache.Cache
}

func NewInvestigationRepoCache(investigationRepo repo.InvestigationRepo, appCache *cache.Cache) *InvestigationRepoCache {
	return &InvestigationRepoCache{
		InvestigationRepo: investigationRepo,
		AppCache:          appCache,
	}
}

func (rc *InvestigationRepoCache) Create(model *model.Investigation) error {
	err := rc.InvestigationRepo.Create(model)
	if err != nil {
		return err
	}

	rc.AppCache.Delete(findAllInvestigationKey(model.ReviewId))
	slog.Debug("InvestigationRepoCache.Create: cache cleared", "reviewId", model.ReviewId)

	return nil
}

func (rc *InvestigationRepoCache) FindAll(reviewID uuid.UUID) (*[]model.Investigation, error) {
	value, found := rc.AppCache.Get(findAllInvestigationKey(reviewID))
	if found {
		slog.Debug("InvestigationRepoCache.FindAll: cache hit", "reviewId", reviewID)
		return value.(*[]model.Investigation), nil
	}
	slog.Debug("InvestigationRepoCache.FindAll: cache miss", "reviewId", reviewID)

	investigations, err := rc.InvestigationRepo.FindAll(reviewID)
	if err != nil {
		return nil, err
	}

	rc.AppCache.Set(findAllInvestigationKey(reviewID), investigations, cache.DefaultExpiration)
	slog.Debug("InvestigationRepoCache.FindAll: cache set", "reviewId", reviewID)

	return investigations, nil
}

func (rc *InvestigationRepoCache) FindOne(investigationId uuid.UUID) (*model.Investigation, error) {
	value, found := rc.AppCache.Get(findOneInvestigationKey(investigationId))
	if found {
		slog.Debug("InvestigationRepoCache.FindOne: cache hit", "investigationId", investigationId)
		return value.(*model.Investigation), nil
	}
	slog.Debug("InvestigationRepoCache.FindOne: cache miss", "investigationId", investigationId)

	investigation, err := rc.InvestigationRepo.FindOne(investigationId)
	if err != nil {
		return nil, err
	}

	rc.AppCache.Set(findOneInvestigationKey(investigationId), investigation, cache.DefaultExpiration)
	slog.Debug("InvestigationRepoCache.FindOne: cache set", "investigationId", investigationId)

	return investigation, nil
}

func (rc *InvestigationRepoCache) SaveKeyword(investigationKeyword *model.InvestigationKeyword) error {
	return rc.InvestigationRepo.SaveKeyword(investigationKeyword)
}

func (rc *InvestigationRepoCache) GetKeywordsByInvestigationId(investigationId uuid.UUID) (*[]model.InvestigationKeyword, error) {
	return rc.InvestigationRepo.GetKeywordsByInvestigationId(investigationId)
}

func (rc *InvestigationRepoCache) GetKeywordsById(keywordId uuid.UUID) (*model.InvestigationKeyword, error) {
	return rc.InvestigationRepo.GetKeywordsById(keywordId)
}

func (rc *InvestigationRepoCache) DeleteKeyword(keywordId uuid.UUID) error {
	return rc.InvestigationRepo.DeleteKeyword(keywordId)
}

func (rc *InvestigationRepoCache) UpdateKeyword(investigationKeyword *model.InvestigationKeyword) (*model.InvestigationKeyword, error) {
	return rc.InvestigationRepo.UpdateKeyword(investigationKeyword)
}
