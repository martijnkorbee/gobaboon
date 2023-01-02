package gobaboon

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/martijnkorbee/gobaboon/internal/pkg/cache"
	"github.com/martijnkorbee/gobaboon/internal/pkg/server"
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
	Cookie server.CookieConfig

	// DatabaseConfig holds the database configuration.
	DatabaseConfig db.DatabaseConfig

	// MailerService sets the mailer service
	MailerService string

	// MailerSettings holds the mailer settings
	MailerSettings mail.MailerSettings
}

// LoadConfig takes in the rootpath of the application
// loads the configuration variables from .env and returns a new config.
// If the .env file is not found or can't be parsed the program will exit.
func mustLoadConfig(path string) Config {
	mustLoadDotEnv(path)

	var (
		debug          bool
		cookieSecure   bool
		cookieLifeTime int
		cookiePersist  bool
		mailerPort     int
	)
	// convert .env variables
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	cookieSecure, err = strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	cookieLifeTime, err = strconv.Atoi(os.Getenv("COOKIE_LIFETIME"))
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	cookiePersist, err = strconv.ParseBool(os.Getenv("COOKIE_PERSIST"))
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	mailerPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	return Config{
		AppName:     os.Getenv("APP_NAME"),
		Rootpath:    path,
		Debug:       debug,
		Host:        os.Getenv("SERVER_HOST"),
		Port:        os.Getenv("SERVER_PORT"),
		RPCport:     os.Getenv("RPC_PORT"),
		Renderer:    os.Getenv("RENDERER"),
		SessionType: os.Getenv("SESSION_TYPE"),
		CacheType:   os.Getenv("CACHE_TYPE"),
		CachePrefix: os.Getenv("CACHE_PREFIX"),
		Redis: cache.RedisConfig{
			Host:     os.Getenv("CACHE_HOST"),
			Port:     os.Getenv("CACHE_PORT"),
			Password: os.Getenv("CACHE_PASS"),
		},
		EncryptionKey: os.Getenv("KEY"),
		Cookie: server.CookieConfig{
			Name:     os.Getenv("COOKIE_NAME"),
			Domain:   os.Getenv("COOKIE_DOMAIN"),
			LifeTime: cookieLifeTime,
			Secure:   cookieSecure,
			Persist:  cookiePersist,
		},
		DatabaseConfig: db.DatabaseConfig{
			Dialect:  os.Getenv("DATABASE_TYPE"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASS"),
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Name:     os.Getenv("DATABASE_NAME"),
			Options: db.Options{
				SSLmode: os.Getenv("DATABASE_SSL_MODE"),
			},
		},
		MailerService: os.Getenv("MAILER"),
		MailerSettings: mail.MailerSettings{
			Domain:     os.Getenv("SMTP_DOMAIN"),
			Host:       os.Getenv("SMTP_HOST"),
			Port:       mailerPort,
			Username:   os.Getenv("SMTP_USERNAME"),
			Password:   os.Getenv("SMTP_PASSWORD"),
			AuthMethod: os.Getenv("SMTP_AUTH_METHOD"),
			Encryption: os.Getenv("SMTP_ENCRYPTION"),
			From:       os.Getenv("SMTP_FROM"),
			FromName:   os.Getenv("SMTP_FROM_NAME"),
		},
	}
}

func mustLoadDotEnv(path string) {
	if _, err := os.Stat(path + "/.env"); os.IsNotExist(err) {
		log.Fatalln("no .env in: ", path)
	}

	err := godotenv.Load(path + "/.env")
	if err != nil {
		log.Fatalln("failed to load .env: ", err)
	}
}
