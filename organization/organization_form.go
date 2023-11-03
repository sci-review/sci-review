package organization

import "golang.org/x/exp/slog"

type OrganizationCreateForm struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description"`
}

func (ocf OrganizationCreateForm) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", ocf.Name),
		slog.String("description", ocf.Description),
	)
}
