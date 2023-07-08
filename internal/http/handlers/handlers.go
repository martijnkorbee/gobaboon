package handlers

import (
	"github.com/martijnkorbee/gobaboon/pkg/logger"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"net/http"
)

type Handlers struct {
	Log    *logger.Logger
	Server *server.Server
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.Server.Renderer.Page(w, r, "home", nil, nil)
	if err != nil {
		h.Log.Error().Err(err).Msg("error rendering")
	}
}
