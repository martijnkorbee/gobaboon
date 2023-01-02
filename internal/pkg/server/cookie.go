package server

import "net/http"

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

// TODO: add cookie methods to set, delete, create new cookies returning http.Cookie
func NewCookie(c CookieConfig) http.Cookie {
	// check default values
	if c.Name == "" {
		c.Name = "baboon"
	}
	if c.Domain == "" {
		c.Domain = "localhost"
	}
	if c.LifeTime == 0 {
		c.LifeTime = 1440
	}
	if c.SameSite == http.SameSiteDefaultMode {
		c.SameSite = http.SameSiteStrictMode
	}

	return http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   c.Secure,
		SameSite: c.SameSite,
		Domain:   c.Domain,
	}
}
