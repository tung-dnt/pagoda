package middleware

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tung-dnt/meme-app/pkg/context"
)

// CacheControl sets a Cache-Control header with a given max age.
func CacheControl(maxAge time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			v := "no-cache, no-store"
			if maxAge > 0 {
				v = fmt.Sprintf("public, max-age=%.0f", maxAge.Seconds())
			}
			ctx.Response().Header().Set("Cache-Control", v)
			return next(ctx)
		}
	}
}

// CacheHTML sets a Cache-Control header suited to edge-caching rendered HTML pages.
//
// A logged-in request is personalized, so it is marked private and never stored by a
// shared cache. An anonymous request is made optimistically cacheable: a short browser
// max-age plus a long stale-while-revalidate window, so Cloudflare can serve the cached
// page instantly and revalidate in the background (via the ETag set by ETag()).
//
// max-age is used deliberately instead of s-maxage: per RFC 9111, s-maxage implies
// proxy-revalidate and would disable stale-while-revalidate at shared caches. The edge
// TTL is controlled separately by a Cloudflare Cache Rule.
func CacheHTML(maxAge, staleWhileRevalidate time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Get(context.AuthenticatedUserKey) != nil {
				ctx.Response().Header().Set("Cache-Control", "private, no-store")
				return next(ctx)
			}

			stale := int(staleWhileRevalidate.Seconds())
			ctx.Response().Header().Set("Cache-Control", fmt.Sprintf(
				"public, max-age=%.0f, stale-while-revalidate=%d, stale-if-error=%d",
				maxAge.Seconds(), stale, stale,
			))
			return next(ctx)
		}
	}
}

// etagBuffer captures a handler's response so a strong ETag can be computed over the
// full body before anything is written to the real client.
type etagBuffer struct {
	http.ResponseWriter
	buf    *bytes.Buffer
	status int
}

func (w *etagBuffer) WriteHeader(status int) { w.status = status }

func (w *etagBuffer) Write(b []byte) (int, error) { return w.buf.Write(b) }

func (w *etagBuffer) Flush() {}

// ETag buffers the response, attaches a strong ETag (sha256 of the body), and answers a
// matching If-None-Match with 304 Not Modified. This is the cheap revalidation path
// Cloudflare uses behind stale-while-revalidate: unchanged content never re-transfers.
//
// It runs inside base's Gzip middleware, so it hashes the uncompressed representation and
// Gzip compresses whatever it writes.
func ETag() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			res := ctx.Response()
			downstream := res.Writer
			capture := &etagBuffer{ResponseWriter: downstream, buf: new(bytes.Buffer), status: http.StatusOK}
			res.Writer = capture

			err := next(ctx)
			res.Writer = downstream
			if err != nil {
				return err
			}

			body := capture.buf.Bytes()
			etag := fmt.Sprintf(`"%x"`, sha256.Sum256(body))
			downstream.Header().Set("ETag", etag)

			if ctx.Request().Header.Get("If-None-Match") == etag {
				downstream.WriteHeader(http.StatusNotModified)
				return nil
			}

			downstream.WriteHeader(capture.status)
			_, werr := downstream.Write(body)
			return werr
		}
	}
}
