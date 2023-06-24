package handlers

import (
	"net/http"

	"github.com/martijnkorbee/gobaboon/pkg/util"
)

// Ping is a default api route
func (h *Handlers) Ping(w http.ResponseWriter, req *http.Request) {
	util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "pong",
	})
}
