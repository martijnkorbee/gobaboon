package app

import (
	"fmt"
	"github.com/martijnkorbee/gobaboon/internal/database/models"
	"github.com/martijnkorbee/gobaboon/internal/http/handlers"
	"github.com/martijnkorbee/gobaboon/internal/http/middleware"
	"github.com/martijnkorbee/gobaboon/internal/http/routes"
	"github.com/martijnkorbee/gobaboon/pkg/db"
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
		app.Log.Fatal().Err(err).Msg("failed load config.properties not allowed")
	}
	app.Config = config

	// create a new server
	app.Server, err = server.NewServer(parseServerConfig(app.Config))
	if err != nil {
		app.Log.Fatal().Err(err).Msg("failed server not allowed")
	}

	// connect to database
	if app.Config.DatabaseConfig.Dialect != "" {
		app.Database, err = mustConnectToDB(&app.Config)
		if err != nil {
			app.Log.Fatal().Err(err).Msg("failed db connection not allowed")
		} else {
			app.Log.Info().Msg("connected to database")
		}
	} else {
		app.Log.Warn().Msg("no database specified")
	}

	// add models
	if app.Config.DatabaseConfig.Dialect != "" {
		app.Models = models.New(app.Database)
	}

	// add middleware
	app.Middleware = &middleware.Middleware{
		Session: app.Server.Session,
		Models:  app.Models,
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

// Path input must be []string{rootpath, databasename}.
func mustConnectToDB(config *Config) (*db.Database, error) {
	var (
		filepath string
	)

	// format sqlite filepath
	if config.DatabaseConfig.Dialect == "sqlite" {
		filepath = fmt.Sprintf("%s/db-data/sqlite/%s.db", config.Rootpath, config.DatabaseConfig.Name)
	}

	// connect to db
	sess, err := db.NewConnection(config.DatabaseConfig, filepath)
	if err != nil {
		return nil, err
	}

	// try connection
	err = sess.Session.Ping()
	if err != nil {
		return nil, err
	}

	return sess, nil
}
