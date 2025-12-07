package teststore

import "github.com/qeery8/rest/internal/app/model"

type TaskRepository struct {
	store *Store
	tasks map[int]string
}

func NewTaskRepository(store *Store) *TaskRepository {
	return &TaskRepository{
		store: store,
		tasks: make(map[int]string),
	}
}

func (r *TaskRepository) Create(t *model.Task) error {
	return nil
}

func (r *TaskRepository) Update(t *model.Task) error {
	return nil
}

func (r *TaskRepository) GetByID(id int) (*model.Task, error) {
	return nil, nil
}

func (r *TaskRepository) Delete(id int) error {
	return nil
}

func (r *TaskRepository) List() ([]*model.Task, error) {
	return nil, nil
}

func (r *TaskRepository) AssigneeUser(userID int, teamID int) error {
	return nil
}
