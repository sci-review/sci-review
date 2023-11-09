package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
	"sci-review/model"
)

type ReviewRepo struct {
	DB    *sqlx.DB
	Cache *cache.Cache
}

func NewReviewRepo(DB *sqlx.DB, cache *cache.Cache) *ReviewRepo {
	return &ReviewRepo{DB: DB, Cache: cache}
}

func findAllKey(userId uuid.UUID) string {
	return "review:findAll:" + userId.String()
}

func findOneKey(id uuid.UUID) string {
	return "review:findOne:" + id.String()
}

func (r *ReviewRepo) Create(review *model.Review, tx *sqlx.Tx) error {
	query := `
		INSERT INTO reviews (id, owner_id, title, type, start_date, end_date, archived, created_at, updated_at)
		VALUES (:id, :owner_id, :title, :type, :start_date, :end_date, :archived, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, review)
	if err != nil {
		return err
	}

	r.Cache.Delete(findAllKey(review.OwnerId))
	slog.Debug("ReviewRepo.Create: cache cleared", "userId", review.OwnerId)

	return nil
}

func (r *ReviewRepo) AddReviewer(reviewer *model.Reviewer, tx *sqlx.Tx) error {
	query := `
		INSERT INTO reviewers (id, user_id, review_id, role, active, created_at, updated_at)
		VALUES (:id, :user_id, :review_id, :role, :active, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, reviewer)
	if err != nil {
		return err
	}

	r.Cache.Delete(findAllKey(reviewer.UserId))
	r.Cache.Delete(findOneKey(reviewer.ReviewId))
	slog.Debug("ReviewRepo.AddReviewer: cache cleared", "userId", reviewer.UserId)

	return nil
}

func (r *ReviewRepo) FindAllByUserId(userId uuid.UUID) (*[]model.Review, error) {
	value, found := r.Cache.Get(findAllKey(userId))
	if found {
		slog.Debug("ReviewRepo.FindAllByUserId: cache hit", "userId", userId)
		return value.(*[]model.Review), nil
	}
	slog.Debug("ReviewRepo.FindAllByUserId: cache miss", "userId", userId)

	var reviews []model.Review
	query := `
		SELECT r.id, r.owner_id, r.title, r.type, r.start_date, r.end_date, r.archived, r.created_at, r.updated_at
		FROM reviews r
		INNER JOIN reviewers rv ON rv.review_id = r.id
		WHERE rv.user_id = $1
	`
	err := r.DB.Select(&reviews, query, userId)
	if err != nil {
		return nil, err
	}

	r.Cache.Set(findAllKey(userId), &reviews, cache.DefaultExpiration)

	return &reviews, nil
}

func (r *ReviewRepo) FindById(id uuid.UUID) (*model.Review, error) {
	value, found := r.Cache.Get(findOneKey(id))
	if found {
		slog.Debug("ReviewRepo.FindById: cache hit", "userId", id)
		return value.(*model.Review), nil
	}
	slog.Debug("ReviewRepo.FindById: cache miss", "userId", id)

	review := model.Review{}
	query := `
		SELECT r.id, r.owner_id, r.title, r.type, r.start_date, r.end_date, r.archived, r.created_at, r.updated_at
		FROM reviews r
		INNER JOIN reviewers rv ON rv.review_id = r.id
		WHERE r.id = $1
	`
	err := r.DB.Get(&review, query, id)
	if err != nil {
		return nil, err
	}

	r.Cache.Set(findOneKey(id), &review, cache.DefaultExpiration)

	return &review, nil
}
