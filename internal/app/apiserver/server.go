package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/qeery8/rest/internal/app/handler"
	"github.com/qeery8/rest/internal/app/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	handlers     handler.Handlers
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.handlers = handler.Handlers{
		User: handler.UserHandlers{
			Store:        store,
			SessionStore: sessionStore,
		},
		Team: handler.TeamHandlers{
			Store: store,
		},
		Task: handler.TaskHandlers{
			Store: store,
		},
	}

	s.configureRouter()

	return s
}
