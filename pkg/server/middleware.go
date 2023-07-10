package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/justinas/nosurf"
	"github.com/martijnkorbee/gobaboon/pkg/util"
)

// Middleware holds all default middleware used by the server.
type Middleware struct {
	rootPath     string
	nosurfCookie http.Cookie
}

// NoSurf is the middleware required to implement CSRF protection.
func (m *Middleware) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// allow api routes
	csrfHandler.ExemptRegexp("/api/.*")

	csrfHandler.SetBaseCookie(m.nosurfCookie)

	return csrfHandler
}

func (m *Middleware) CheckMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if MaintenanceMode {

			allowedURLS := strings.Split(os.Getenv("ALLOWED_URLS"), ",")
			allowedURLS = append(allowedURLS, "public/static/html/maintenance.html")

			allowed := func() bool {
				for _, url := range allowedURLS {
					if strings.Contains(r.URL.Path, url) {
						return true
					}
				}
				return false
			}()

			if !allowed {
				if strings.Contains(r.URL.Path, "/api/") {
					_ = util.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
						"message": "under maintenance",
					})
					return
				}

				w.WriteHeader(http.StatusServiceUnavailable)
				w.Header().Set("Retry-After:", "300")
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				http.ServeFile(w, r, fmt.Sprintf("%s/public/static/html/maintenance.html", m.rootPath))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
