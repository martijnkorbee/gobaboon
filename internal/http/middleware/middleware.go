package middleware

import (
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/pkg/server"
)

type Middleware struct {
	Session *server.Session
	Models  *models.Models
}
