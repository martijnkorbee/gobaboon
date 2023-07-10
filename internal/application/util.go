package application

import (
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"net/http"
)

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
