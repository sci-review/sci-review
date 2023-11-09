package form

import (
	"golang.org/x/exp/slog"
)

type InvestigationForm struct {
	Question string `json:"question" form:"question" validate:"required,min=3"`
}

func (i InvestigationForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("question", i.Question),
	)
}
