package sqlstore

import (
	"time"

	"github.com/qeery8/rest/internal/app/model"
)

type TeamRepository struct {
	store *Store
}

func (r *TeamRepository) Create(t *model.Team) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	if err := r.store.db.QueryRow(
		`INSERT INTO teams (name, description, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id `,
		t.Name, t.Description, t.OwnerID, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID); err != nil {
		return err
	}

	r.AddMembers(t.OwnerID, t.ID)

	return nil
}

func (r *TeamRepository) AddMembers(userID int, teamID int) error {
	_, err := r.store.db.Exec(
		`INSERT INTO team_members (user_id, team_id, created_at)
		VALUES ($1, $2, NOW())`,
		userID, teamID,
	)
	return err
}

func (r *TeamRepository) Find(id int) (*model.Team, error) {
	t := &model.Team{}
	err := r.store.db.QueryRow(
		"SELECT id, name, description, owner_id, created_at, updated_at FROM teams WHERE id = $1",
		id,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.OwnerID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TeamRepository) FindByUser(userID int) ([]*model.Team, error) {
	rows, err := r.store.db.Query(
		`SELECT t.id, t.name, t.description, t.owner_id, t.created_at, t.updated_at
		FROM teams t
		JOIN team_members tm ON tm.team_id = t.id
		WHERE tm.user_id = $1`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teams []*model.Team

	for rows.Next() {
		t := &model.Team{}
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.OwnerID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) RemoveMembers(userID int, teamID int) error {
	_, err := r.store.db.Exec(
		`DELETE FROM team_members
		WHERE user_id = $1 AND team_id = $2`,
		userID, teamID,
	)
	return err
}

func (r *TeamRepository) Update(t *model.Team) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.UpdatedAt = time.Now()

	_, err := r.store.db.Exec(
		`UPDATE teams SET
		name = $1, description = $2, updated_at = $3
		WHERE id = $4`,
		t.Name, t.Description, t.UpdatedAt, t.ID,
	)

	return err
}

func (r *TeamRepository) Delete(id int) error {
	_, err := r.store.db.Exec(
		`DELETE FROM teams
		WHERE id = $1`,
		id,
	)
	return err
}
