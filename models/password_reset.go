package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	// MinBytesPerToken - the min number of bytes to be used for each session token.
	// MinBytesPerToken     = 32
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less then the
	// MinBytesPerToken const it will be ignored
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (ps *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement password reeset service")
}

func (ps *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement password reset consume ")
}
