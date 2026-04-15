package model

import "time"

type TaskStatus string

type TaskPriority string

const (
	StatusToDo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

const (
	LowPriority    TaskPriority = "low"
	MediumPriority TaskPriority = "medium"
	HightPriority  TaskPriority = "hight"
)

type Task struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Content    string       `json:"content"`
	Status     TaskStatus   `json:"status"`
	Priority   TaskPriority `json:"priority"`
	DueDate    *time.Time   `json:"due_date"`
	AssigneeID *int         `json:"assignee_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func (t *Task) Validate() error {
	if len(t.Name) < 4 {
		return ErrNameShort
	}
	if len(t.Name) > 50 {
		return ErrNameTooLong
	}
	return nil
}
