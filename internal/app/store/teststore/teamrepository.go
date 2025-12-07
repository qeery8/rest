package teststore

import "github.com/qeery8/rest/internal/app/model"

type TeamRepository struct {
	store        *Store
	teams        map[int]*model.Team
	team_members map[int]map[int]bool
}

func NewTeamRepository(store *Store) *TeamRepository {
	return &TeamRepository{
		store:        store,
		teams:        make(map[int]*model.Team),
		team_members: make(map[int]map[int]bool),
	}
}

func (r *TeamRepository) Create(t *model.Team) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return nil
}

func (r *TeamRepository) Update(t *model.Team) error {
	return nil
}

func (r *TeamRepository) AddMembers(userID int, teamID int) error {
	return nil
}

func (r *TeamRepository) RemoveMembers(userID int, teamID int) error {
	return nil
}

func (r *TeamRepository) Find(id int) (*model.Team, error) {
	return nil, nil
}

func (r *TeamRepository) FindByUser(userID int) ([]*model.Team, error) {
	return nil, nil
}

func (r *TeamRepository) Delete(id int) error {
	return nil
}
