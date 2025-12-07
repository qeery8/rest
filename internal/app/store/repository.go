package store

import "github.com/qeery8/rest/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Update(*model.User) error
	Delete(id int) error
}

type TeamRepository interface {
	Create(*model.Team) error
	Find(int) (*model.Team, error)
	FindByUser(userID int) ([]*model.Team, error)
	Update(*model.Team) error
	Delete(id int) error
	AddMembers(teamID int, userID int) error
	RemoveMembers(teamID int, userID int) error
}

type TaskRepository interface {
	Create(*model.Task) error
	AssigneeUser(userID int, taskID int) error
	Update(*model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	List() ([]*model.Task, error)
}
