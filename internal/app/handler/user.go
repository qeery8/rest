package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/qeery8/rest/internal/app/ctxkeys"
	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/session"
	"github.com/qeery8/rest/internal/app/store"
	"github.com/qeery8/rest/internal/app/utils"
)

type UserHandlers struct {
	Store        store.Store
	SessionStore sessions.Store
}

func (s *UserHandlers) HandleUsersCreate() http.HandlerFunc {
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

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.Store.User().Create(u); err != nil {
			utils.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		utils.Respond(w, r, http.StatusCreated, u)
	}
}

func (s *UserHandlers) HandlerWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Respond(w, r, http.StatusOK, r.Context().Value(ctxkeys.CtxKeyUser).(*model.User))
	}

}

func (s *UserHandlers) HandleUpdateProfile() http.HandlerFunc {
	type request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value(ctxkeys.CtxKeyUser).(*model.User)

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		if !currentUser.ComparePassword(req.OldPassword) {
			utils.Error(w, r, http.StatusUnauthorized, errors.New("ivalid old password"))
			return
		}

		currentUser.Password = req.NewPassword
		if err := s.Store.User().Update(currentUser); err != nil {
			utils.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, map[string]string{
			"message": "password was updated",
		})
	}
}

func (s *UserHandlers) HandlerUsersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.SessionStore.Get(r, session.SessionsName)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values = make(map[interface{}]interface{})
		session.Options.MaxAge = -1

		if err := s.SessionStore.Save(r, w, session); err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *UserHandlers) HandlerDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value(ctxkeys.CtxKeyUser).(*model.User)

		if err := s.Store.User().Delete(currentUser.ID); err != nil {
			utils.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		sessions, _ := s.SessionStore.Get(r, session.SessionsName)
		delete(sessions.Values, "user_id")
		s.SessionStore.Save(r, w, sessions)

		utils.Respond(w, r, http.StatusOK, map[string]string{
			"message": "account deleted successfully",
		})
	}
}
