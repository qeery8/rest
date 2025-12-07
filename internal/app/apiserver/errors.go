package apiserver

import "errors"

var (
	ErrTeamNotFound  = errors.New("team not found")
	ErrInvalidTeamId = errors.New("invalid team id")
	ErrTaskNotFound  = errors.New("task not found")
)
