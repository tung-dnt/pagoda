package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/tung-dnt/meme-app/pkg/context"
	"github.com/tung-dnt/meme-app/pkg/tests"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCacheHTML_Anonymous(t *testing.T) {
	// An anonymous request (no authenticated user in context) is edge-cacheable:
	// max-age for the browser, a long stale-while-revalidate window, and stale-if-error.
	ctx, _ := tests.NewContext(c.Web, "/about")

	_ = tests.ExecuteMiddleware(ctx, CacheHTML(60*time.Second, 24*time.Hour))

	assert.Equal(t,
		"public, max-age=60, stale-while-revalidate=86400, stale-if-error=86400",
		ctx.Response().Header().Get("Cache-Control"))
}

func TestCacheHTML_Authenticated(t *testing.T) {
	// A logged-in request must never be shared-cached: the page is personalized.
	ctx, _ := tests.NewContext(c.Web, "/about")
	ctx.Set(context.AuthenticatedUserKey, usr)

	_ = tests.ExecuteMiddleware(ctx, CacheHTML(60*time.Second, 24*time.Hour))

	assert.Equal(t, "private, no-store", ctx.Response().Header().Get("Cache-Control"))
}

func TestETag_SetsStrongETagAndServesBody(t *testing.T) {
	// With no validator on the request, the handler runs and a strong ETag of the
	// exact body is attached to a normal 200 response.
	ctx, rec := tests.NewContext(c.Web, "/about")
	body := "<html>hello world</html>"

	err := tests.ExecuteHandler(ctx, func(c echo.Context) error {
		return c.HTML(http.StatusOK, body)
	}, ETag())
	assert.NoError(t, err)

	want := fmt.Sprintf(`"%x"`, sha256.Sum256([]byte(body)))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, want, rec.Header().Get("ETag"))
	assert.Equal(t, body, rec.Body.String())
}

func TestETag_NotModifiedWhenValidatorMatches(t *testing.T) {
	// A matching If-None-Match short-circuits to 304 with an empty body — the cheap
	// revalidation path Cloudflare uses behind stale-while-revalidate.
	ctx, rec := tests.NewContext(c.Web, "/about")
	body := "<html>hello world</html>"
	etag := fmt.Sprintf(`"%x"`, sha256.Sum256([]byte(body)))
	ctx.Request().Header.Set("If-None-Match", etag)

	err := tests.ExecuteHandler(ctx, func(c echo.Context) error {
		return c.HTML(http.StatusOK, body)
	}, ETag())
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotModified, rec.Code)
	assert.Empty(t, rec.Body.String())
	assert.Equal(t, etag, rec.Header().Get("ETag"))
}

func TestETag_ServesBodyWhenValidatorDiffers(t *testing.T) {
	// A stale validator must not short-circuit: the fresh body is served with the new ETag.
	ctx, rec := tests.NewContext(c.Web, "/about")
	body := "<html>fresh content</html>"
	ctx.Request().Header.Set("If-None-Match", `"deadbeef"`)

	err := tests.ExecuteHandler(ctx, func(c echo.Context) error {
		return c.HTML(http.StatusOK, body)
	}, ETag())
	assert.NoError(t, err)

	want := fmt.Sprintf(`"%x"`, sha256.Sum256([]byte(body)))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, want, rec.Header().Get("ETag"))
	assert.Equal(t, body, rec.Body.String())
}
