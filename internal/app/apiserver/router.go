package apiserver

import (
	"net/http"

	gh "github.com/gorilla/handlers"
)

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(gh.CORS(gh.AllowedOrigins([]string{"*"})))
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/users", s.handlers.User.HandleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/teams", s.handlers.Team.HandleTeamsCreate()).Methods("POST")
	s.router.HandleFunc("/task", s.handlers.Task.HandleTaskCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)

	private.HandleFunc("/whoami", s.handlers.User.HandlerWhoami()).Methods("GET")
	private.HandleFunc("/team/{id}", s.handlers.Team.HandleTeamID()).Methods("GET")
	private.HandleFunc("/team_user_id/{id}", s.handlers.Team.HandleTeamByUserID()).Methods("GET")
	private.HandleFunc("/task/list", s.handlers.Task.HandleTaskList()).Methods("GET")
	private.HandleFunc("/task/{task_id}", s.handlers.Task.HandleTaskGetID()).Methods("GET")

	private.HandleFunc("/teams/{team_id}/members", s.handlers.Team.HandleTeamAddMembers()).Methods("POST")
	private.HandleFunc("/task/{user_id}/member", s.handlers.Task.HandleTaskAssigneeID()).Methods("POST")
	private.HandleFunc("/logout", s.handlers.User.HandlerUsersDelete()).Methods("POST")

	private.HandleFunc("/profile", s.handlers.User.HandleUpdateProfile()).Methods("PUT")
	private.HandleFunc("/team/{team_id}", s.handlers.Team.HandleTeamUpdate()).Methods("PUT")
	private.HandleFunc("/task/{task_id}", s.handlers.Task.HandleTaskUpdate()).Methods("PUT")

	private.HandleFunc("/user/me", s.handlers.User.HandlerDelete()).Methods("DELETE")
	private.HandleFunc("/task/{task_id}", s.handlers.Task.HandleTaskDelete()).Methods("DELETE")
	private.HandleFunc("/team/{team_id}/members/{user_id}", s.handlers.Team.HandleTeamMembersDelete()).Methods("DELETE")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
