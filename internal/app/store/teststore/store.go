package teststore

import (
	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
)

// store ...
type Store struct {
	userRepository *UserRepository
	teamRepository *TeamRepository
	taskRepository *TaskRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}

func (s *Store) Team() store.TeamRepository {
	if s.teamRepository != nil {
		return s.teamRepository
	}

	s.teamRepository = &TeamRepository{
		store:        s,
		teams:        make(map[int]*model.Team),
		team_members: make(map[int]map[int]bool),
	}

	return s.teamRepository
}

func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
		tasks: make(map[int]string),
	}

	return s.taskRepository
}
