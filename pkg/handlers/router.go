package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/tung-dnt/meme-app/pkg/context"
	"github.com/tung-dnt/meme-app/pkg/middleware"
	"github.com/tung-dnt/meme-app/pkg/services"
	files "github.com/tung-dnt/meme-app/public"
)

// BuildRouter builds the router.
func BuildRouter(c *services.Container) error {
	// Force HTTPS, if enabled.
	if c.Config.HTTP.TLS.Enabled {
		c.Web.Use(echomw.HTTPSRedirect())
	}

	// Serve public files with cache control.
	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.PublicFile)).
		Static("files", "public/files")

	// Serve static files.
	// ui.StaticFile() should be used in ui components to append a cache key to the URL to break cache
	// after each server reboot.
	c.Web.Group(
		"",
		echomw.GzipWithConfig(echomw.GzipConfig{
			Skipper: func(c echo.Context) bool {
				for _, ext := range []string{
					".js",
					".css",
				} {
					if strings.HasSuffix(c.Request().URL.Path, ext) {
						return false
					}
				}
				return true
			},
		}),
		middleware.CacheControl(c.Config.Cache.Expiration.PublicFile),
	).StaticFS("static", echo.MustSubFS(files.Static, "static"))

	// Create a cookie store for session data.
	cookieStore := sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.SameSite = http.SameSiteStrictMode

	// Base middleware shared by both dynamic route groups. Defined once so the authenticated
	// and public groups never repeat it. Note Session and LoadAuthenticatedUser are here (not
	// CSRF): they only read cookies, so public pages stay cookie-free for anonymous visitors
	// while still rendering personalized content for logged-in ones.
	base := []echo.MiddlewareFunc{
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		middleware.SetLogger(),
		middleware.LogRequest(),
		echomw.Gzip(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: c.Config.App.Timeout,
		}),
		middleware.Config(c.Config),
		middleware.Session(cookieStore),
		middleware.LoadAuthenticatedUser(c.Auth),
	}

	// Authenticated group: base + CSRF. Hosts every form/action route.
	g := c.Web.Group("", base...)
	g.Use(echomw.CSRFWithConfig(echomw.CSRFConfig{
		TokenLookup:    "form:csrf",
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
		ContextKey:     context.CSRFKey,
	}))

	// Public group: base + optimistic HTML caching and ETag revalidation, and deliberately no
	// CSRF (which would set a per-session cookie and fragment the shared cache). Only
	// cookie-free, form-less pages may be registered here.
	p := c.Web.Group("", base...)
	p.Use(
		middleware.CacheHTML(c.Config.Cache.Expiration.PublicPage, c.Config.Cache.StaleWhileRevalidate),
		middleware.ETag(),
	)

	// Error handler.
	c.Web.HTTPErrorHandler = new(Error).Page

	// Initialize and register all handlers.
	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			return err
		}

		h.Routes(g, p)
	}

	return nil
}
