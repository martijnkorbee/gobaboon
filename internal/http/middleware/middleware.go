package middleware

import (
	"baboonapp/database/models"

	"github.com/martijnkorbee/gobaboon"
)

type Middleware struct {
	App    *gobaboon.Baboon
	Models *models.Models
}
