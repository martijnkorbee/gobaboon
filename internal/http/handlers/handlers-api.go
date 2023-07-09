package handlers

import (
	"github.com/martijnkorbee/gobaboon/pkg/util"
	"net/http"
)

// Ping is a default api route (used to check health)
func (h *Handlers) Ping(w http.ResponseWriter, req *http.Request) {
	_ = util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "pong",
	})
}
