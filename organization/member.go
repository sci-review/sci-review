package organization

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"time"
)

type Member struct {
	Id             uuid.UUID  `db:"id" json:"id"`
	UserId         uuid.UUID  `db:"user_id" json:"userId"`
	OrganizationId uuid.UUID  `db:"organization_id" json:"organizationId"`
	Role           MemberRole `db:"role" json:"role"`
	Active         bool       `db:"active" json:"active"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}

func NewMember(userId uuid.UUID, organizationId uuid.UUID, role MemberRole, active bool) *Member {
	return &Member{
		Id:             uuid.New(),
		UserId:         userId,
		OrganizationId: organizationId,
		Role:           role,
		Active:         active,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (m Member) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", m.Id.String()),
		slog.String("user_id", m.UserId.String()),
		slog.String("organization_id", m.OrganizationId.String()),
		slog.String("role", string(m.Role)),
		slog.Bool("active", m.Active),
		slog.Time("created_at", m.CreatedAt),
		slog.Time("updated_at", m.UpdatedAt),
	)
}
