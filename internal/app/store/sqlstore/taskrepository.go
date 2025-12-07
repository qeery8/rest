package sqlstore

import (
	"time"

	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Create(t *model.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	return r.store.db.QueryRow(
		`INSERT INTO tasks (name, content, status, priority, due_date, assignee_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		t.Name, t.Content, t.Status, t.Priority, t.DueDate, t.AssigneeID, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID)
}

func (r *TaskRepository) Update(t *model.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.UpdatedAt = time.Now()

	_, err := r.store.db.Exec(
		`UPDATE tasks SET 
		name = $1, content = $2, status = $3, priority = $4, due_date = $5, assignee_id = $6, updated_at = $7
		WHERE id = $8`,
		t.Name, t.Content, t.Status, t.Priority, t.DueDate, t.AssigneeID, t.UpdatedAt, t.ID,
	)

	return err
}

func (r *TaskRepository) GetByID(id int) (*model.Task, error) {
	t := &model.Task{}
	err := r.store.db.QueryRow(
		`SELECT id, name, content, status, priority, due_date, assignee_id, created_at, updated_at
		FROM tasks 
		WHERE id = $1`, id,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Content,
		&t.Status,
		&t.Priority,
		&t.DueDate,
		&t.AssigneeID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TaskRepository) List() ([]*model.Task, error) {
	rows, err := r.store.db.Query(
		`SELECT id, name, content, status, priority, due_date, assignee_id, created_at, updated_at
		FROM tasks`,
	)

	if err != nil {
		return nil, err
	}

	var tasks []*model.Task

	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Content,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.AssigneeID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) Delete(id int) error {
	_, err := r.store.db.Exec(
		`DELETE FROM tasks
		WHERE id = $1`, id,
	)
	return err
}

func (r *TaskRepository) AssigneeUser(userID int, teamID int) error {
	query := `
		UPDATE tasks 
		SET assignee_id = $1 
		WHERE id = $2 
		AND EXISTS (
			SELECT 1 
			FROM team_members
			WHERE team_members.user_id = $1
			AND team_members.team_id = tasks.team_id
		)
	`
	result, err := r.store.db.Exec(query, userID, teamID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return store.ErrUserNotInTeam
	}
	return nil
}
