package model

import (
	"github.com/google/uuid"
)

type Principal struct {
	Id   uuid.UUID `json:"id"`
	Role Role      `json:"role"`
}

func NewPrincipal(id string, role string) *Principal {
	return &Principal{
		Id:   uuid.MustParse(id),
		Role: Role(role),
	}
}
