package gobaboon

import (
	"net/http"

	"github.com/martijnkorbee/gobaboon/internal/pkg/cache"
	"github.com/martijnkorbee/gobaboon/pkg/db"
	"github.com/martijnkorbee/gobaboon/pkg/mail"
)

// Config holds all configuration settings to be used throughout baboon.
type Config struct {
	AppName string

	// Rootpath is the rootpath of the application. Usually this the full path to the bin file.
	Rootpath string

	// Debug is used to set some functionality in debug mode. Mostly being more explicit in log info and errors.
	Debug bool

	// Host will be used by the server.
	Host string

	// Port will be used by the server.
	Port string

	// RPC port if set baboon also starts an rpc listener
	RPCport string

	// EncryptionKey is used to encrypt and decrypt with their respective functions.
	EncryptionKey string

	// Renderer sets which type of template engine will be used.
	Renderer string

	// SessionType sets which type of session store to use i.e. cookie, cache, db.
	SessionType string

	// CacheType sets which cache client to use
	CacheType string

	// CachePrefix sets the cache prefix for the server
	CachePrefix string

	// RedisConfig holds redis client configuration
	Redis cache.RedisConfig

	// Cookie holds cookie configuration
	Cookie CookieConfig

	// DatabaseConfig holds the database configuration.
	DatabaseConfig db.DatabaseConfig

	// MailerService sets the mailer service
	MailerService string

	// MailerSettings holds the mailer settings
	MailerSettings mail.MailerSettings
}

type CookieConfig struct {
	// Name defaults to baboon
	Name string

	// Domain defaults to localhost
	Domain string

	// Lifetime defaults to 1440 minutes
	LifeTime int // time in minutes

	// Secure defaults to false
	Secure bool

	// Persist defaults to false
	Persist bool

	// SameSite defaults to SameSiteStrict mode
	SameSite http.SameSite
}
