package apiserver

import (
	"net/http"

	gh "github.com/gorilla/handlers"
)

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(gh.CORS(gh.AllowedOrigins([]string{"*"})))
	s.router.Use(s.logRequest)

	//регистрация пользователя
	s.router.HandleFunc("/users", s.handlers.User.HandleUsersCreate()).Methods("POST")
	//авторизация
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	//создание команды
	s.router.HandleFunc("/teams", s.handlers.Team.HandleTeamsCreate()).Methods("POST")
	//создание задачи
	s.router.HandleFunc("/task", s.handlers.Task.HandleTaskCreate()).Methods("POST")

	//это типо приватные запросы, хуй знает как объяснить
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)

	//для проверки авторизованного пользователя, те по запросу выдает инфу из бд о челе
	private.HandleFunc("/whoami", s.handlers.User.HandlerWhoami()).Methods("GET")
	//выдает инфу о команде по введенному id
	private.HandleFunc("/team/{id}", s.handlers.Team.HandleTeamID()).Methods("GET")
	//выдает инфу о командах в которых состоит юзер
	private.HandleFunc("/team_user_id/{id}", s.handlers.Team.HandleTeamByUserID()).Methods("GET")
	//показывает все задачи которые тебе присвоены
	private.HandleFunc("/task/list", s.handlers.Task.HandleTaskList()).Methods("GET")
	//показывает задачу по id
	private.HandleFunc("/task/{task_id}", s.handlers.Task.HandleTaskGetID()).Methods("GET")

	//добавляет юзера в команду
	private.HandleFunc("/teams/{team_id}/members", s.handlers.Team.HandleTeamAddMembers()).Methods("POST")
	//присваивает задачу юзеру
	private.HandleFunc("/task/{user_id}/member", s.handlers.Task.HandleTaskAssigneeID()).Methods("POST")
	//заканчивает активную сессию
	private.HandleFunc("/logout", s.handlers.User.HandlerUsersDelete()).Methods("POST")

	//обновляет пароль (не помню делал ли чтобы можно было почту поменять или нет)
	private.HandleFunc("/profile", s.handlers.User.HandleUpdateProfile()).Methods("PUT")
	//обновляет название, описание и тд команды
	private.HandleFunc("/team/{team_id}", s.handlers.Team.HandleTeamUpdate()).Methods("PUT")
	//обновляет название, контент задачи и тд (короче если что потом просто уточнишь)
	private.HandleFunc("/task/{task_id}", s.handlers.Task.HandleTaskUpdate()).Methods("PUT")

	//удаляет пользователя из бд (свой акк)
	private.HandleFunc("/user/me", s.handlers.User.HandlerDelete()).Methods("DELETE")
	//удаляет задачу по id
	private.HandleFunc("/task/{task_id}", s.handlers.Task.HandleTaskDelete()).Methods("DELETE")
	//удаляет участника команды
	private.HandleFunc("/team/{team_id}/members/{user_id}", s.handlers.Team.HandleTeamMembersDelete()).Methods("DELETE")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
