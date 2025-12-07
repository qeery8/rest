package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/qeery8/rest/internal/app/ctxkeys"
	"github.com/qeery8/rest/internal/app/errors"
	"github.com/qeery8/rest/internal/app/session"
	"github.com/qeery8/rest/internal/app/utils"
	"github.com/sirupsen/logrus"
)

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request_ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxkeys.CtxKeyRequsetID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxkeys.CtxKeyRequsetID),
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
		session, err := s.sessionStore.Get(r, session.SessionsName)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			utils.Error(w, r, http.StatusUnauthorized, errors.ErrNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			utils.Error(w, r, http.StatusUnauthorized, errors.ErrNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxkeys.CtxKeyUser, u)))
	})
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			utils.Error(w, r, http.StatusUnauthorized, errors.ErrIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, session.SessionsName)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, nil)
	}
}
