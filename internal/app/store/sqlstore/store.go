package sqlstore

import (
	"database/sql"

	"github.com/qeery8/rest/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	teamRepository *TeamRepository
	taskRepository *TaskRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Team() store.TeamRepository {
	if s.teamRepository != nil {
		return s.teamRepository
	}

	s.teamRepository = &TeamRepository{
		store: s,
	}

	return s.teamRepository
}

func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
	}

	return s.taskRepository
}
