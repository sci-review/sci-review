package organization

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"time"
)

type Organization struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func NewOrganization(name string, description string) *Organization {
	return &Organization{
		Id:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (o Organization) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", o.Id.String()),
		slog.String("name", o.Name),
		slog.String("description", o.Description),
		slog.Time("created_at", o.CreatedAt),
		slog.Time("updated_at", o.UpdatedAt),
	)
}
