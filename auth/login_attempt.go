package auth

import (
	"github.com/google/uuid"
	"time"
)

type LoginAttempt struct {
	Id        uuid.UUID     `db:"id" json:"id"`
	UserID    uuid.NullUUID `db:"user_id" json:"userId"`
	Email     string        `db:"email" json:"email"`
	Success   bool          `db:"success" json:"success"`
	IPAddress string        `db:"ip_address" json:"ipAddress"`
	UserAgent string        `db:"user_agent" json:"userAgent"`
	Timestamp time.Time     `db:"timestamp" json:"timestamp"`
}

func NewLoginAttempt(userID uuid.NullUUID, email string, success bool, IPAddress string, userAgent string) *LoginAttempt {
	return &LoginAttempt{
		Id:        uuid.New(),
		UserID:    userID,
		Email:     email,
		Success:   success,
		IPAddress: IPAddress,
		UserAgent: userAgent,
		Timestamp: time.Now(),
	}
}

func NewUnSuccessLoginAttempt(email string, IPAddress string, userAgent string) *LoginAttempt {
	return NewLoginAttempt(uuid.NullUUID{}, email, false, IPAddress, userAgent)
}

func NewSuccessLoginAttempt(userID uuid.UUID, email string, IPAddress string, userAgent string) *LoginAttempt {
	return NewLoginAttempt(uuid.NullUUID{UUID: userID, Valid: true}, email, true, IPAddress, userAgent)
}
