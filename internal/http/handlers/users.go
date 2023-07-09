package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/pkg/util"
	"net/http"
	"time"
)

func (h *Handlers) UsersAdd(w http.ResponseWriter, req *http.Request) {
	user := models.User{
		Active: 1,
	}

	// read request
	err := util.ReadJSON(w, req, &user)
	if err != nil {
		h.Log.Error().Err(err).Msg("error reading from request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := h.Models.Users.AddUser(user)
	if err != nil {
		h.Log.Error().Err(err).Msg("error adding user")
		_ = util.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	token, err := h.Models.Tokens.GenerateToken(id, 720*time.Hour)
	if err != nil {
		h.Log.Error().Err(err).Msg("error generating token")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	usr, err := h.Models.Users.GetUserByID(id)
	if err != nil {
		h.Log.Error().Err(err).Msg("error retrieving user from database")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.Models.Tokens.InsertToken(*token, *usr)
	if err != nil {
		h.Log.Error().Err(err).Msg("error inserting user token")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_ = util.WriteJSON(w, http.StatusOK, map[string]string{
		"email": user.Email,
		"token": token.PlainText,
	})
}

func (h *Handlers) UsersDelete(w http.ResponseWriter, req *http.Request) {
	user := models.User{
		Active: 1,
	}

	email := chi.URLParam(req, "email")

	usr, err := h.Models.Users.GetUserByEmail(email)
	if err != nil {
		h.Log.Error().Err(err).Msg("error retrieving user")
		_ = util.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = h.Models.Users.DeleteUser(usr.ID)
	if err != nil {
		h.Log.Error().Err(err).Msg("error deleting user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// delete all tokens for user
	for {
		tkn, err := h.Models.Tokens.GetTokenByEmail(user.Email)
		if err != nil {
			break
		}
		_ = h.Models.Tokens.DeleteTokenByID(tkn.ID)
	}

	_ = util.WriteJSON(w, http.StatusOK, map[string]string{
		"action": "delete",
		"user":   user.Email,
	})
}
