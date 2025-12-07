package model

import (
	"errors"
	"time"
)

var (
	ErrNameShort   = errors.New("team name so short")
	ErrNameTooLong = errors.New("team name too long")
)

type Team struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	OwnerID     int       `json:"owner_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Team) Validate() error {
	if len(t.Name) < 3 {
		return ErrNameShort
	}
	if len(t.Name) > 20 {
		return ErrNameTooLong
	}
	return nil
}
