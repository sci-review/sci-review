package auth

import (
	"github.com/google/uuid"
	"sci-review/user"
)

type Principal struct {
	Id   uuid.UUID `json:"id"`
	Role user.Role `json:"role"`
}

func NewPrincipal(id string, role string) *Principal {
	return &Principal{
		Id:   uuid.MustParse(id),
		Role: user.Role(role),
	}
}
