package gobaboon

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgraph-io/badger"
	"github.com/gomodule/redigo/redis"
	"github.com/martijnkorbee/gobaboon/internal/pkg/cache"
	"github.com/martijnkorbee/gobaboon/internal/pkg/rpc"
	"github.com/martijnkorbee/gobaboon/internal/pkg/server"
	"github.com/martijnkorbee/gobaboon/pkg/db"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
	"github.com/martijnkorbee/gobaboon/pkg/mail"
	"github.com/martijnkorbee/gobaboon/pkg/util"
	"github.com/robfig/cron/v3"
)

// const version = "1.0.0"

type Baboon struct {
	// Config holds all baboon required config settings
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

	// Mailer is the baboon app mailer
	Mailer *mail.Mailer

	// Database holds baboon's main database
	Database *db.Database

	// Cache is baboon's cache client
	Cache cache.Cache
}

// New creates a new baboon app
func (b *Baboon) Init(rootpath string) error {
	// load config
	b.Config = mustLoadConfig(rootpath)

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

	// check if not exist create expected folders for a new baboon app
	fnames := []string{
		"http/handlers",
		"http/middleware",
		"database/models",
		"database/migrations",
		"database/data",
		"templates/views",
		"templates/mail",
		"logs",
	}
	for _, fname := range fnames {
		err := util.CreateDirIfNotExists(rootpath + "/" + fname)
		if err != nil {
			return err
		}
	}

	// create scheduler
	b.Scheduler = cron.New()

	// create RPC server
	b.RPCServer = rpc.NewRPCServer(b.Config.RPCport)

	// connect to database
	if b.Config.DatabaseConfig.Dialect != "" {
		b.Database = b.mustConnectToDB()
	}

	// connect to cache
	if b.Config.CacheType != "" {
		b.Cache = b.mustConnectToCache()
	}

	// create a new server
	if srv, err := server.NewServer(parseServerConfig(b.Config)); err != nil {
		return err
	} else {
		b.Server = srv
	}

	// if session type is set to a persistent store
	if b.Config.SessionType != "" && b.Config.SessionType != "cookie" {
		switch b.Config.SessionType {
		case "postgers", "postgresql", "sqlite", "mysql":
			b.Server.Session.SetPersistentStoreDatabase(b.Database.Dialect, b.Database.Session.Driver().(*sql.DB))
		case "redis":
			b.Server.Session.SetPersistentStoreRedis(b.Cache.GetConnection().(*redis.Pool))
		case "badger":
			b.Server.Session.SetPersistentStoreBadger(b.Cache.GetConnection().(*badger.DB))
		default:
			b.Log.Fatal().Err(errors.New("unsported session store")).Msg("failed to set persistent session store")
		}
	}

	// create a new mailer
	if b.Config.MailerService != "" {
		if mail, err := mail.NewMailer(parseMailConfig(b.Config), b.Config.MailerService, rootpath+"/templates/mail"); err != nil {
			return err
		} else {
			b.Mailer = mail
		}
	}

	return nil
}

func (b *Baboon) Run() error {
	// close database session when server terminates
	if b.Config.DatabaseConfig.Dialect != "" {
		defer b.Database.Session.Close()
	}

	// close cache connection when server terminates
	if b.Config.CacheType != "" {
		defer b.Cache.Close()
	}

	// add badger garbage collection to scheduler
	if b.Config.CacheType == "badger" {
		cid, err := b.Scheduler.AddFunc("@daily", func() { b.Cache.GetConnection().(*badger.DB).RunValueLogGC(0.7) })
		if err != nil {
			b.Log.Error().Err(err).Msg("failed to a GC for badger")
		} else {
			b.Log.Info().Str("cronEntryID", fmt.Sprint(cid)).Msg("added GC for badger to cron jobs")
		}
	}

	// start mail listener
	if b.Config.MailerService != "" {
		b.Log.Info().Msg("starting mail channels")
		go b.Mailer.ListenForMail()
	}

	// start RPC server
	b.Log.Info().Str("port", b.RPCServer.Port).Msg("starting RPC server")
	go b.RPCServer.Run()

	// start web server
	b.Log.Info().Str("port", b.Config.Port).Msg("starting web server")
	if err := b.Server.Run(); err != nil {
		b.Log.Fatal().Err(err).Msg("failed to start web server")
	}

	return nil
}

// parseServerConfig helps to extract baboon configation variables into a server configuration
func parseServerConfig(c Config) server.ServerConfig {
	return server.ServerConfig{
		Rootpath: c.Rootpath,
		Host:     c.Host,
		Port:     c.Port,
		Renderer: c.Renderer,
		Debug:    c.Debug,
		Cookie: server.CookieConfig{
			Name:     c.Cookie.Name,
			Domain:   c.Cookie.Domain,
			LifeTime: c.Cookie.LifeTime,
			Persist:  c.Cookie.Persist,
			Secure:   c.Cookie.Secure,
			SameSite: http.SameSiteStrictMode,
		},
	}
}

func parseMailConfig(c Config) mail.MailerSettings {
	return mail.MailerSettings{
		Domain:     c.MailerSettings.Domain,
		Host:       c.MailerSettings.Host,
		Port:       c.MailerSettings.Port,
		Username:   c.MailerSettings.Username,
		Password:   c.MailerSettings.Password,
		AuthMethod: c.MailerSettings.AuthMethod,
		Encryption: c.MailerSettings.Encryption,
		From:       c.MailerSettings.From,
		FromName:   c.MailerSettings.FromName,
	}
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
		filepath = fmt.Sprintf("%s/database/data/sqlite/%s.db", b.Config.Rootpath, b.Config.DatabaseConfig.Name)
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
			Path:   b.Config.Rootpath + "/database/data/badger",
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
