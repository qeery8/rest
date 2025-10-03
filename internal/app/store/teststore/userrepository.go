package teststore

import (
	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func NewUserRepository(store *Store) *UserRepository {
	return &UserRepository{
		store: store,
		users: make(map[int]*model.User),
	}
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1

	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) Update(u *model.User) error {
	targetUser, ok := r.users[u.ID]
	if !ok {
		return store.ErrRecordNotFound
	}

	if u.Email != "" {
		targetUser.Email = u.Email
	}

	if u.Password != "" {
		if err := u.BeforeCreate(); err != nil {
			return err
		}
		targetUser.EncryptedPassword = u.EncryptedPassword
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	if _, ok := r.users[id]; !ok {
		return store.ErrRecordNotFound
	}

	delete(r.users, id)

	return nil
}
