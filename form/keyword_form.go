package form

import "golang.org/x/exp/slog"

type KeywordForm struct {
	Word     string `json:"word" form:"word" validate:"required,min=3,max=255"`
	Synonyms string `json:"synonyms" form:"synonyms"`
}

func (k KeywordForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("word", k.Word),
		slog.String("synonyms", k.Synonyms),
	)
}
