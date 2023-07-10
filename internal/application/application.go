package application

import (
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/internal/http/handlers"
	"github.com/martijnkorbee/gobaboon/internal/http/middleware"
	"github.com/martijnkorbee/gobaboon/internal/http/routes"
	"github.com/martijnkorbee/gobaboon/pkg/cache"
	"github.com/martijnkorbee/gobaboon/pkg/db"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"github.com/robfig/cron/v3"
	"sync"
)

type Application struct {
	// Config holds all required .config.properties settings
	Config Config

	// WG is used to block before exiting
	WG sync.WaitGroup

	// Log is the application's default logger
	Log *logger.Logger

	// Scheduler can be used to schedule cron jobs
	Scheduler *cron.Cron

	// Server
	Server     *server.Server
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Routes     *routes.AppRoutes

	Database *db.Database
	Models   *models.Models

	Cache *cache.Cache
}
