package model

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"time"
)

type Organization struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	Archived    bool      `db:"archived" json:"archived"`
	Members     []Member  `db:"-" json:"members"`
}

func NewOrganization(name string, description string) *Organization {
	return &Organization{
		Id:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Archived:    false,
		Members:     []Member{},
	}
}

func (o *Organization) AddMember(member Member) {
	o.Members = append(o.Members, member)
}

func (o Organization) IsActiveMember(userId uuid.UUID) bool {
	for _, member := range o.Members {
		if member.UserId == userId && member.Active {
			return true
		}
	}
	return false
}

func (o Organization) IsOwner(userId uuid.UUID) bool {
	for _, member := range o.Members {
		if member.UserId == userId && member.Role == MemberOwner {
			return true
		}
	}
	return false
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
