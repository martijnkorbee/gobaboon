package application

import (
	"github.com/joho/godotenv"
	"github.com/martijnkorbee/gobaboon/pkg/cache"
	"net/http"
	"os"
	"strconv"

	"github.com/martijnkorbee/gobaboon/pkg/db"
)

// Config holds all configuration settings.
type Config struct {
	AppName string

	// Rootpath is the rootpath of the Application. Usually this the full path to the bin file.
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

// LoadConfig takes in the rootpath of the application
// loads the config.properties and returns a the config.
func mustLoadConfig(path string) (Config, error) {
	var (
		debug bool
	)

	if _, err := os.Stat(path + "/.config.properties"); os.IsNotExist(err) {
		return Config{}, err
	}

	// load config
	if err := godotenv.Load(path + "/.config.properties"); err != nil {
		return Config{}, err
	}

	// set default values
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = true
	}

	if os.Getenv("SERVER_HOST") == "" {
		if err := os.Setenv("SERVER_HOST", "0.0.0.0"); err != nil {
			return Config{}, err
		}
	}

	if os.Getenv("SERVER_PORT") == "" {
		if err := os.Setenv("SERVER_PORT", "4000"); err != nil {
			return Config{}, err
		}
	}

	if os.Getenv("SESSION_TYPE") == "" {
		if err := os.Setenv("SESSION_TYPE", "cookie"); err != nil {
			return Config{}, err
		}
	}

	return Config{
		AppName:       os.Getenv("APP_NAME"),
		Rootpath:      path,
		Debug:         debug,
		Host:          os.Getenv("SERVER_HOST"),
		Port:          os.Getenv("SERVER_PORT"),
		Renderer:      os.Getenv("RENDERER"),
		SessionType:   os.Getenv("SESSION_TYPE"),
		CacheType:     os.Getenv("CACHE_TYPE"),
		CachePrefix:   os.Getenv("CACHE_PREFIX"),
		EncryptionKey: os.Getenv("KEY"),
		Cookie: CookieConfig{
			Name:   os.Getenv("COOKIE_NAME"),
			Domain: os.Getenv("COOKIE_DOMAIN"),
		},
		DatabaseConfig: db.DatabaseConfig{
			Dialect: os.Getenv("DATABASE_TYPE"),
			Name:    os.Getenv("DATABASE_NAME"),
		},
	}, nil
}
