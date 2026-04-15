package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrUserNotInTeam  = errors.New("user not in team")
)
