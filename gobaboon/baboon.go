package gobaboon

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/gomodule/redigo/redis"
	"github.com/martijnkorbee/gobaboon/pkg/cache"
	"github.com/martijnkorbee/gobaboon/pkg/db"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
	"github.com/martijnkorbee/gobaboon/pkg/rpc"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"github.com/robfig/cron/v3"
)

// const version = "1.0.0"

type Baboon struct {
	// Config holds all baboon required .config.properties settings
	Config Config

	// Log is the default logger for baboon
	// Applications using baboon should assign their own loggers as in the skeleton app.
	Log *logger.Logger

	// Scheduler can be used to schedule tasks (like cron jobs)
	Scheduler *cron.Cron

	// Server is the baboon app server.
	Server *server.Server

	// RPCServer is baboon's RPC server
	RPCServer *rpc.RPCServer

	// Database holds baboon's main database
	Database *db.Database

	// Cache is baboon's cache client
	Cache cache.Cache
}

// New creates a new baboon app
func (b *Baboon) Init(c Config) error {
	// set .config.properties
	b.Config = c

	// start baboon logger
	log := &logger.LoggerConfig{
		Rootpath:   b.Config.Rootpath,
		Debug:      b.Config.Debug,
		Console:    b.Config.Debug,
		ToFile:     !b.Config.Debug,
		Service:    "baboon-main",
		Filename:   "/logs/baboon_log.log",
		MaxBackups: 2,
		LocalTime:  true,
	}
	b.Log = log.Start()

	// create scheduler
	b.Scheduler = cron.New()

	// create RPC server
	// TODO: disabled
	//b.RPCServer = rpc.NewRPCServer(b.Config.Host, b.Config.RPCport)

	// connect to database
	if b.Config.DatabaseConfig.Dialect != "" {
		b.Database = b.mustConnectToDB()
	}

	// connect to cache
	if b.Config.CacheType != "" {
		b.Cache = b.mustConnectToCache()
	}

	// if session type is set to a persistent store
	if b.Config.SessionType != "" && b.Config.SessionType != "cookie" {
		switch b.Config.SessionType {
		case "postgress", "postgresql", "sqlite", "mysql":
			if err := b.Server.Session.SetPersistentStoreDatabase(b.Database.Dialect, b.Database.Session.Driver().(*sql.DB)); err != nil {
				//	do something
			}
		case "redis":
			b.Server.Session.SetPersistentStoreRedis(b.Cache.GetConnection().(*redis.Pool))
		case "badger":
			b.Server.Session.SetPersistentStoreBadger(b.Cache.GetConnection().(*badger.DB))
		default:
			b.Log.Fatal().Err(errors.New("unsported session store")).Msg("failed to set persistent session store")
		}
	}

	return nil
}

// mustConnectToDB must connect to database, if db type is sqlite formats the connection.
//
// Path input must be []string{rootpath, databasename}.
func (b *Baboon) mustConnectToDB() *db.Database {
	var (
		filepath string
	)

	// format sqlite filepath
	if b.Config.DatabaseConfig.Dialect == "sqlite" {
		filepath = fmt.Sprintf("%s/db-models/sqlite/%s.db", b.Config.Rootpath, b.Config.DatabaseConfig.Name)
	}

	// connect to db
	sess, err := db.NewConnection(b.Config.DatabaseConfig, filepath)
	if err != nil {
		b.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// try connection
	err = sess.Session.Ping()
	if err != nil {
		b.Log.Fatal().Err(err).Msg("failed to ping database")
	}

	b.Log.Info().Msg(fmt.Sprintf("connected to database: %s", b.Config.DatabaseConfig.Dialect))

	return sess
}

// mustConnectToCache must connect to cache
func (b *Baboon) mustConnectToCache() cache.Cache {
	switch b.Config.CacheType {
	case "redis":
		client, err := cache.CreateRedisCache(cache.RedisConfig{
			Prefix:   b.Config.CachePrefix,
			Host:     b.Config.Redis.Host,
			Port:     b.Config.Redis.Port,
			Password: b.Config.Redis.Password,
		})
		if err != nil {
			b.Log.Fatal().Err(err).Msg("failed to connect to redis")
		}
		// assign redis client
		return client
	case "badger":
		client, err := cache.CreateBadgerCache(cache.BadgerConfig{
			Prefix: b.Config.CachePrefix,
			Path:   b.Config.Rootpath + "/database/models/badger",
		}, b.Log)
		if err != nil {
			b.Log.Fatal().Err(err).Msg("failed to connect to redis")
		}
		// assign badger client
		return client
	default:
		b.Log.Fatal().Err(errors.New("unsupported client type")).Msg("failed to connect to cache")
		return nil
	}
}
