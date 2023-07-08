package app

import (
	"fmt"
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/internal/http/handlers"
	"github.com/martijnkorbee/gobaboon/internal/http/middleware"
	"github.com/martijnkorbee/gobaboon/internal/http/routes"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/martijnkorbee/gobaboon"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
)

type application struct {
	Baboon     *gobaboon.Baboon
	Log        *logger.Logger
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Routes     *routes.AppRoutes
	Models     *models.Models
	WG         sync.WaitGroup
}

// Starts starts the application
func (a *application) Start() {
	// add your things to do before starting here

	// spawn shutdown listener
	go a.listenForShutdown()

	// start baboon
	a.Log.Info().Msg("starting baboon")
	if err := a.Baboon.Run(); err != nil {
		a.Log.Fatal().Err(err).Msg("failed to start baboon")
	}
}

// listenForShutDown captures shutdown signals and initiates shutdown process
func (a *application) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	a.Log.Info().Msg(fmt.Sprintf("received signal: %s", s.String()))

	// initiate shutdown
	a.shutdown()
}

// shutdown gracefully shuts down the application
func (a *application) shutdown() {
	// put any clean up tasks here

	// block untill the wait group is empty
	a.WG.Wait()

	// exit application
	os.Exit(0)
}
