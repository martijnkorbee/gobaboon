package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/badgerstore"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
	"github.com/gomodule/redigo/redis"
)

type Session struct {
	Manager *scs.SessionManager
	server  *Server
}

// New creates and returns a new session manager. Needs a cookieconfig.
func NewSession(c CookieConfig, s *Server) *Session {
	var (
		sess = scs.New()
	)

	// set session options
	sess.Cookie.Name = c.Name
	sess.Lifetime = time.Duration(c.LifeTime) * time.Minute
	sess.Cookie.Persist = c.Persist
	sess.Cookie.Secure = c.Secure
	sess.Cookie.Domain = c.Domain
	sess.Cookie.SameSite = http.SameSiteStrictMode

	return &Session{
		Manager: sess,
		server:  s,
	}
}

// SetPersistentStoreDatabase sets the session store to the database, expects required session schema to be present.
//
// The bobo CLI command <bobo make session> can create and execute the corresponding db migration for you.
// Below example is for postgres, refer to https://github.com/alexedwards/scs#configuring-the-session-store for full documentation.
/*
CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	models BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
*/
func (s *Session) SetPersistentStoreDatabase(dialect string, driver *sql.DB) error {
	switch strings.ToLower(dialect) {
	case "postgres", "postgresql":
		s.Manager.Store = postgresstore.New(driver)
	case "mysql":
		s.Manager.Store = mysqlstore.New(driver)
	case "sqlite":
		s.Manager.Store = sqlite3store.New(driver)
	default:
		return errors.New("failed to set persistent store, db type not supported")
	}

	return nil
}

// SetPersistenStoreRedis sets the session store to redis
func (s *Session) SetPersistentStoreRedis(pool *redis.Pool) {
	s.Manager.Store = redisstore.New(pool)
}

// SetPersistentStoreBadger sets the session store to badger
func (s *Session) SetPersistentStoreBadger(badger *badger.DB) {
	s.Manager.Store = badgerstore.New(badger)
}
