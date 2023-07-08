package app

import (
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/internal/http/handlers"
	"github.com/martijnkorbee/gobaboon/internal/http/middleware"
	"github.com/martijnkorbee/gobaboon/internal/http/routes"
	"log"
	"os"

	"github.com/martijnkorbee/gobaboon"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
)

func MustInitApplication() *application {
	// get application root path
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	// init baboon
	baboon := &gobaboon.Baboon{}
	err = baboon.Init(gobaboon.Config{
		Rootpath:    path,
		Debug:       true,
		Host:        "localhost",
		Port:        "4000",
		Renderer:    "jet",
		SessionType: "cookie",
		// LOAD YOUR CONFIGURATION HERE YOU COULD USE THE .env EXAMPLE
	})
	if err != nil {
		log.Fatalln(err)
	}
	// AFTER ADDING YOUR CONFIG YOU CAN REMOVE THIS LOG
	baboon.Log.Warn().Msg("no config loaded, running with default settings")

	app := &application{
		Baboon: baboon,
	}

	// start application logger
	log := &logger.LoggerConfig{
		Rootpath:   app.Baboon.Config.Rootpath,
		Debug:      app.Baboon.Config.Debug,
		Console:    app.Baboon.Config.Debug,
		ToFile:     !app.Baboon.Config.Debug,
		Service:    "application-main",
		Filename:   "/logs/app_log.log",
		MaxBackups: 2,
		LocalTime:  true,
	}
	app.Log = log.Start()

	// add models
	if app.Baboon.Config.DatabaseConfig.Dialect != "" {
		app.Models = models.New(app.Baboon.Database)
	}

	// add middleware
	app.Middleware = &middleware.Middleware{
		App:    baboon,
		Models: app.Models,
	}

	// add handlers
	app.Handlers = &handlers.Handlers{
		App: baboon,
	}

	// add routes
	app.Routes = &routes.AppRoutes{
		Middleware: *app.Middleware,
		Handlers:   *app.Handlers,
	}

	// mount application routes
	app.Baboon.Server.Router.Mount("/", app.Routes.Routes())       // app routes
	app.Baboon.Server.Router.Mount("/api", app.Routes.RoutesAPI()) // API routes

	return app
}
