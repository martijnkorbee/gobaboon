package gobaboon

import (
	"database/sql"
	"errors"
	"github.com/dgraph-io/badger"
	"github.com/gomodule/redigo/redis"
	"github.com/martijnkorbee/gobaboon/pkg/cache"
	"github.com/martijnkorbee/gobaboon/pkg/rpc"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"github.com/robfig/cron/v3"
)

// const version = "1.0.0"

type Baboon struct {
	// Config holds all baboon required .config.properties settings
	Config Config

	// Scheduler can be used to schedule tasks (like cron jobs)
	Scheduler *cron.Cron

	// Server is the baboon app server.
	Server *server.Server

	// RPCServer is baboon's RPC server
	RPCServer *rpc.RPCServer

	// Cache is baboon's cache client
	Cache cache.Cache
}

// New creates a new baboon app
func (b *Baboon) Init(c Config) error {
	// set .config.properties
	b.Config = c

	// create scheduler
	b.Scheduler = cron.New()

	// create RPC server
	// TODO: disabled
	//b.RPCServer = rpc.NewRPCServer(b.Config.Host, b.Config.RPCport)

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
