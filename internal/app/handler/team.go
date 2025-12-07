package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qeery8/rest/internal/app/errors"
	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
	"github.com/qeery8/rest/internal/app/utils"
)

type TeamHandlers struct {
	Store store.Store
}

func (s *TeamHandlers) HandleTeamsCreate() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		OwnerID     int    `json:"owner_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		t := &model.Team{
			Name:        req.Name,
			Description: req.Description,
			OwnerID:     req.OwnerID,
		}

		if err := s.Store.Team().Create(t); err != nil {
			utils.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *TeamHandlers) HandleTeamID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idStr, ok := vars["id"]
		if !ok {
			utils.Error(w, r, http.StatusBadRequest, errors.ErrTaskNotFound)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, errors.ErrInvalidTeamId)
			return
		}

		team, err := s.Store.Team().Find(id)
		if err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, team)
	}
}

func (s *TeamHandlers) HandleTeamByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idStr, ok := vars["id"]
		if !ok {
			utils.Error(w, r, http.StatusBadRequest, errors.ErrTeamNotFound)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, errors.ErrInvalidTeamId)
			return
		}

		teams, err := s.Store.Team().FindByUser(id)
		if err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, teams)
	}
}

func (s *TeamHandlers) HandleTeamAddMembers() http.HandlerFunc {
	type request struct {
		UserID int `json:"user_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		teamIDStr := vars["team_id"]

		teamID, err := strconv.Atoi(teamIDStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.Store.Team().AddMembers(teamID, req.UserID); err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}
		utils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *TeamHandlers) HandleTeamMembersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		teamIdStr := vars["team_id"]
		userIdStr := vars["user_id"]

		teamID, _ := strconv.Atoi(teamIdStr)
		userID, _ := strconv.Atoi(userIdStr)

		if err := s.Store.Team().RemoveMembers(teamID, userID); err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}
		utils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *TeamHandlers) HandleTeamUpdate() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["team_id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, errors.ErrInvalidTeamId)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		team := &model.Team{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := s.Store.Team().Update(team); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		updates, err := s.Store.Team().Find(id)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, updates)
	}
}
