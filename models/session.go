package models

import (
	"database/sql"
	"fmt"
	"github.com/mihailtudos/photosharer/rand"
)

type Session struct {
	ID     int
	UserID int

	//Token is only set when creating a new session. When look up a session
	// this will be left empty, as we only store the hash of the session
	// in our DB and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss SessionService) Create(userID int) (*Session, error) {
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID: userID,
		Token:  token,
	}

	return &session, nil
}

func (ss SessionService) User(token string) (*User, error) {
	return nil, nil
}
