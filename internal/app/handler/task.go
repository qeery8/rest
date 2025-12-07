package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
	"github.com/qeery8/rest/internal/app/utils"
)

type TaskHandlers struct {
	Store store.Store
}

func (s *TaskHandlers) HandleTaskCreate() http.HandlerFunc {
	type request struct {
		Name     string             `json:"name"`
		Content  string             `json:"content"`
		Status   model.TaskStatus   `json:"status"`
		Priority model.TaskPriority `json:"priority"`
		DueDate  *time.Time         `json:"due_date"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		t := &model.Task{
			Name:     req.Name,
			Content:  req.Content,
			Status:   req.Status,
			Priority: req.Priority,
			DueDate:  req.DueDate,
		}

		if err := s.Store.Task().Create(t); err != nil {
			utils.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		utils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *TaskHandlers) HandleTaskList() http.HandlerFunc {
	type request struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		tasks, err := s.Store.Task().List()
		if err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}
		utils.Respond(w, r, http.StatusOK, tasks)
	}
}

func (s *TaskHandlers) HandleTaskUpdate() http.HandlerFunc {
	type request struct {
		Name       string             `json:"name"`
		Content    string             `json:"content"`
		Status     model.TaskStatus   `json:"status"`
		Priority   model.TaskPriority `json:"priority"`
		AssigneeID *int               `json:"assignee_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["task_id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		task := &model.Task{
			ID:         id,
			Name:       req.Name,
			Content:    req.Content,
			Status:     req.Status,
			Priority:   req.Priority,
			AssigneeID: req.AssigneeID,
		}

		if err := s.Store.Task().Update(task); err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}

		updated, err := s.Store.Task().GetByID(id)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, updated)
	}
}

func (s *TaskHandlers) HandleTaskGetID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["task_id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		task, err := s.Store.Task().GetByID(id)
		if err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}
		utils.Respond(w, r, http.StatusOK, task)
	}
}

func (s *TaskHandlers) HandleTaskDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["task_id"]
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.Store.Task().Delete(taskID); err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *TaskHandlers) HandleTaskAssigneeID() http.HandlerFunc {
	type request struct {
		TeamID int `json:"team_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["user_id"]
		userID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := s.Store.Task().AssigneeUser(userID, req.TeamID); err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.Respond(w, r, http.StatusOK, nil)
	}
}
