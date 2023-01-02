// Package db facades the creation of different db connections and adapter imports in order to return a generalised db type.
// This makes it easier to create a new db connection throughout the application without having to worry about imports.
// Full documentatation: https://upper.io/v4/getting-started/.
package db

import (
	"errors"
	"fmt"

	"github.com/martijnkorbee/gobaboon/pkg/util"
	upper "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mssql"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
	"github.com/upper/db/v4/adapter/sqlite"
)

type Database struct {
	Dialect string
	Session upper.Session
}

type DatabaseConfig struct {
	Dialect  string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Options  Options
}

type Options struct {
	SSLmode string
}

// NewConnection takes in a database configuration, for sqlite a single full file path,
// returns a new database connection layer.
func NewConnection(c DatabaseConfig, path ...string) (*Database, error) {
	var (
		curl upper.ConnectionURL
	)

	switch c.Dialect {
	case "postgres", "postgresql":
		curl = postgresql.ConnectionURL{
			User:     c.User,
			Password: c.Password,
			Host:     c.Host,
			Socket:   c.Port,
			Database: c.Name,
			Options: map[string]string{
				"sslmode": c.Options.SSLmode,
			},
		}
	case "mysql", "mariadb":
		url, err := mysql.ParseURL(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			c.User,
			c.Password,
			c.Host,
			c.Port,
			c.Name,
		))
		if err != nil {
			return nil, err
		}

		url.Options = map[string]string{
			"tls": c.Options.SSLmode,
		}
		curl = url
	case "sqlite":
		if len(path) == 0 {
			return nil, errors.New("no full file path given for sqlite db")
		}

		if len(path) > 1 {
			return nil, errors.New("to many arguments given for sqlite db")
		}

		// create db if it doesn't exist
		err := util.CreateFileIfNotExists(path[0])
		if err != nil {
			return nil, err
		}

		curl = sqlite.ConnectionURL{
			Database: path[0],
			Options:  map[string]string{},
		}
	case "mssql", "sqlserver":
		curl = mssql.ConnectionURL{
			User:     c.User,
			Password: c.Password,
			Database: c.Name,
			Host:     c.Host,
			Socket:   c.Port,
			Options:  map[string]string{},
		}
	default:
		return nil, errors.New("non supported database type")
	}

	sess, err := upper.Open(c.Dialect, curl)
	if err != nil {
		return nil, err
	}

	return &Database{
		Session: sess,
	}, nil
}
