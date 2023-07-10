package application

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/martijnkorbee/gobaboon/pkg/cache"
	"github.com/rs/zerolog/log"
	upper "github.com/upper/db/v4"
	"os"
	"os/signal"
	"syscall"
)

// Start starts the Application
func (a *Application) Init() {
	// add your things to do before starting here

	// spawn shutdown listener
	go a.listenForShutdown()
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

func (a *Application) Run() error {
	// close database session when the application stops
	if a.Config.DatabaseConfig.Dialect != "" {
		defer func(session upper.Session) {
			if err := session.Close(); err != nil {
				log.Err(err).Timestamp().Stack().Msg("error closing database connection")
			}
		}(a.Database.Session)
	}

	// close cache connection when server terminates
	if a.Config.CacheType != "" {
		defer func(cache cache.Cache) {
			if err := cache.Close(); err != nil {
				log.Err(err).Timestamp().Stack().Msg("error closing cache connection")
			}
		}(*a.Cache)
	}

	// add badger garbage collection to scheduler
	if a.Config.CacheType == "badger" {
		cid, err := a.Scheduler.AddFunc("@daily", func() { a.Cache.GetConnection().(*badger.DB).RunValueLogGC(0.7) })
		if err != nil {
			a.Log.Error().Err(err).Msg("failed to a GC for badger")
		} else {
			a.Log.Info().Str("cronEntryID", fmt.Sprint(cid)).Msg("added GC for badger to cron jobs")
		}
	}

	// start app server
	a.Log.Info().Str("port", a.Config.Port).Msg("starting app server")
	if err := a.Server.Run(); err != nil {
		return err
	}

	return nil
}
