package form

import (
	"golang.org/x/exp/slog"
	"sci-review/model"
)

type ReviewCreateForm struct {
	Title      string           `json:"title" form:"title" validate:"required,min=3,max=255"`
	ReviewType model.ReviewType `json:"review_type" form:"review_type" validate:"required,oneof=SystematicReview ScopingReview RapidReview"`
	StartDate  string           `json:"start_date" form:"start_date" validate:"required"`
	EndDate    string           `json:"end_date" form:"end_date" validate:"required"`
}

func (r ReviewCreateForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", r.Title),
		slog.String("review_type", string(r.ReviewType)),
		slog.String("start_date", r.StartDate),
		slog.String("end_date", r.EndDate),
	)
}
