package app

import (
	"github.com/martijnkorbee/gobaboon/internal/http/handlers"
	"github.com/martijnkorbee/gobaboon/internal/http/middleware"
	"github.com/martijnkorbee/gobaboon/internal/http/routes"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"log"
	"os"

	"github.com/martijnkorbee/gobaboon/pkg/logger"
)

func New() *Application {
	// get Application root path
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	app := &Application{}

	// start Application logger
	// TODO: settings as arguments
	applog := &logger.LoggerConfig{
		Rootpath:  app.Config.Rootpath,
		Debug:     app.Config.Debug,
		Console:   true,
		Service:   "Application",
		LocalTime: true,
	}
	app.Log = applog.Start()

	// #######################################################################
	// you can add your own way of loading env or use this as a starting point
	// #######################################################################

	// load configuration
	config, err := mustLoadConfig(path)
	if err != nil {
		app.Log.Fatal().Err(err).Msg("failed to load config.properties")
	}
	app.Config = config

	// create a new server
	if srv, err := server.NewServer(parseServerConfig(app.Config)); err != nil {
		app.Log.Fatal().Err(err).Msg("failed to create new server")
	} else {
		app.Server = srv
	}

	//// add models
	//if app.Config.DatabaseConfig.Dialect != "" {
	//	app.Models = models.New(app.Database)
	//}

	// add middleware
	app.Middleware = &middleware.Middleware{
		Models: app.Models,
	}

	// add handlers
	app.Handlers = &handlers.Handlers{
		Server: app.Server,
		Log:    app.Log,
	}

	// add routes
	app.Routes = &routes.AppRoutes{
		Middleware: *app.Middleware,
		Handlers:   *app.Handlers,
	}

	// mount Application routes
	app.Server.Router.Mount("/", app.Routes.Routes())       // app routes
	app.Server.Router.Mount("/api", app.Routes.RoutesAPI()) // API routes

	return app
}
