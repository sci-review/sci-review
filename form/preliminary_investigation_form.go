package form

import (
	"golang.org/x/exp/slog"
)

type PreliminaryInvestigationForm struct {
	Question string `json:"question" form:"question" validate:"required,min=3"`
}

func (pi PreliminaryInvestigationForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("question", pi.Question),
	)
}
