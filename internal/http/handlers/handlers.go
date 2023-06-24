package handlers

import (
	"net/http"

	"github.com/martijnkorbee/gobaboon"
)

type Handlers struct {
	App *gobaboon.Baboon
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.App.Server.Renderer.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.Log.Error().Err(err).Msg("error rendering")
	}
}
