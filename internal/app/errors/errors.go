package errors

import "errors"

var (
	ErrTeamNotFound  = errors.New("team not found")
	ErrInvalidTeamId = errors.New("invalid team id")
	ErrTaskNotFound  = errors.New("task not found")
)

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrNotAuthenticated         = errors.New("not authenticated")
)
