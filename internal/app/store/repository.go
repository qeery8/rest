package store

import "github.com/qeery8/rest/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Update(*model.User) error
	Delete(id int) error
}
