package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "Kolyan_pidr"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequsetID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/teams", s.handleTeamsCreate()).Methods("POST")
	s.router.HandleFunc("/task", s.handleTaskCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)

	private.HandleFunc("/whoami", s.handlerWhoami()).Methods("GET")
	private.HandleFunc("/team/{id}", s.handleTeamID()).Methods("GET")
	private.HandleFunc("/team_user_id/{id}", s.handleTeamByUserID()).Methods("GET")
	private.HandleFunc("/task/list", s.handleTaskList()).Methods("GET")
	private.HandleFunc("/task/{task_id}", s.handleTaskGetID()).Methods("GET")

	private.HandleFunc("/teams/{team_id}/members", s.handleTeamAddMembers()).Methods("POST")
	private.HandleFunc("/task/{user_id}/member", s.handleTaskAssigneeID()).Methods("POST")
	private.HandleFunc("/logout", s.handlerUsersDelete()).Methods("POST")

	private.HandleFunc("/profile", s.handleUpdateProfile()).Methods("PUT")
	private.HandleFunc("/team/{team_id}", s.handleTeamUpdate()).Methods("PUT")
	private.HandleFunc("/task/{task_id}", s.handleTaskUpdate()).Methods("PUT")

	private.HandleFunc("/user/me", s.handlerDelete()).Methods("DELETE")
	private.HandleFunc("/task/{task_id}", s.handleTaskDelete()).Methods("DELETE")
	private.HandleFunc("/team/{team_id}/members/{user_id}", s.handleTeamMembersDelete()).Methods("DELETE")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request_ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequsetID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequsetID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
			code:           http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
		)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handlerWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}

}

func (s *server) handleUpdateProfile() http.HandlerFunc {
	type request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value(ctxKeyUser).(*model.User)

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if !currentUser.ComparePassword(req.OldPassword) {
			s.error(w, r, http.StatusUnauthorized, errors.New("ivalid old password"))
			return
		}

		currentUser.Password = req.NewPassword
		if err := s.store.User().Update(currentUser); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{
			"message": "password was updated",
		})
	}
}

func (s *server) handlerUsersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values = make(map[interface{}]interface{})
		session.Options.MaxAge = -1

		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handlerDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value(ctxKeyUser).(*model.User)

		if err := s.store.User().Delete(currentUser.ID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		sessions, _ := s.sessionStore.Get(r, sessionName)
		delete(sessions.Values, "user_id")
		s.sessionStore.Save(r, w, sessions)

		s.respond(w, r, http.StatusOK, map[string]string{
			"message": "account deleted successfully",
		})
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTeamsCreate() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		OwnerID     int    `json:"owner_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		t := &model.Team{
			Name:        req.Name,
			Description: req.Description,
			OwnerID:     req.OwnerID,
		}

		if err := s.store.Team().Create(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTeamID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idStr, ok := vars["id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, ErrTeamNotFound)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidTeamId)
			return
		}

		team, err := s.store.Team().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, team)
	}
}

func (s *server) handleTeamByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idStr, ok := vars["id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, ErrTeamNotFound)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidTeamId)
			return
		}

		teams, err := s.store.Team().FindByUser(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, teams)
	}
}

func (s *server) handleTeamAddMembers() http.HandlerFunc {
	type request struct {
		UserID int `json:"user_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		teamIDStr := vars["team_id"]

		teamID, err := strconv.Atoi(teamIDStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Team().AddMembers(teamID, req.UserID); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTeamMembersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		teamIdStr := vars["team_id"]
		userIdStr := vars["user_id"]

		teamID, _ := strconv.Atoi(teamIdStr)
		userID, _ := strconv.Atoi(userIdStr)

		if err := s.store.Team().RemoveMembers(teamID, userID); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTeamUpdate() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["team_id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidTeamId)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		team := &model.Team{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := s.store.Team().Update(team); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		updates, err := s.store.Team().Find(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, updates)
	}
}

func (s *server) handleTaskCreate() http.HandlerFunc {
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		t := &model.Task{
			Name:     req.Name,
			Content:  req.Content,
			Status:   req.Status,
			Priority: req.Priority,
			DueDate:  req.DueDate,
		}

		if err := s.store.Task().Create(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTaskList() http.HandlerFunc {
	type request struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		tasks, err := s.store.Task().List()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, tasks)
	}
}

func (s *server) handleTaskUpdate() http.HandlerFunc {
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
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

		if err := s.store.Task().Update(task); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		updated, err := s.store.Task().GetByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, updated)
	}
}

func (s *server) handleTaskGetID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["task_id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		task, err := s.store.Task().GetByID(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, task)
	}
}

func (s *server) handleTaskDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["task_id"]
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Task().Delete(taskID); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTaskAssigneeID() http.HandlerFunc {
	type request struct {
		TeamID int `json:"team_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["user_id"]
		userID, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := s.store.Task().AssigneeUser(userID, req.TeamID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
