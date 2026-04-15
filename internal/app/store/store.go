package store

type Store interface {
	User() UserRepository
	Team() TeamRepository
	Task() TaskRepository
}
