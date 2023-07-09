package app

import (
	"fmt"
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/internal/http/handlers"
	"github.com/martijnkorbee/gobaboon/internal/http/middleware"
	"github.com/martijnkorbee/gobaboon/internal/http/routes"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/martijnkorbee/gobaboon/pkg/logger"
)

type Application struct {
	Config Config
	Server *server.Server
	// TODO: add db
	//Database   *db.Database
	// TODO: add cache
	Log        *logger.Logger
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Routes     *routes.AppRoutes
	Models     *models.Models
	WG         sync.WaitGroup
}

// Start starts the Application
func (a *Application) Start() {
	// add your things to do before starting here

	// spawn shutdown listener
	go a.listenForShutdown()

	// start baboon
	a.Log.Info().Msg("starting app")
	if err := a.Run(); err != nil {
		a.Log.Fatal().Err(err).Msg("failed to start app")
	}
}

// listenForShutDown captures shutdown signals and initiates shutdown process
func (a *Application) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	a.Log.Info().Msg(fmt.Sprintf("received signal: %s", s.String()))

	// initiate shutdown
	a.shutdown()
}

// shutdown gracefully shuts down the Application
func (a *Application) shutdown() {
	// put any clean up tasks here

	// block until the wait group is empty
	a.WG.Wait()

	// exit Application
	os.Exit(0)
}
