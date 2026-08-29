package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/tung-dnt/meme-app/pkg/routenames"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Public application pages must be edge-cacheable: an anonymous GET carries the optimistic
// Cache-Control, a strong ETag, and — crucially — no CSRF cookie (which would fragment the
// shared cache and block it entirely).
func TestPublicRoutes_AreCacheableAndCookieless(t *testing.T) {
	for _, route := range []string{routenames.Home, routenames.About} {
		resp := request(t).setRoute(route).get().assertStatusCode(http.StatusOK)

		cc := resp.Header.Get("Cache-Control")
		assert.Contains(t, cc, "public", "route %q should be publicly cacheable", route)
		assert.Contains(t, cc, "stale-while-revalidate=", "route %q should allow stale serving", route)
		assert.NotEmpty(t, resp.Header.Get("ETag"), "route %q should carry a strong ETag", route)

		for _, ck := range resp.Header.Values("Set-Cookie") {
			assert.NotContains(t, ck, "_csrf", "public route %q must not set the CSRF cookie", route)
		}
	}
}

// A matching If-None-Match short-circuits to 304 — the full stale-while-revalidate pipeline.
func TestPublicRoutes_ConditionalGetReturns304(t *testing.T) {
	url := srv.URL + c.Web.Reverse(routenames.About)

	first, err := http.Get(url)
	require.NoError(t, err)
	etag := first.Header.Get("ETag")
	require.NotEmpty(t, etag)
	first.Body.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("If-None-Match", etag)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotModified, resp.StatusCode)
}

// Regression guard: form pages stay in the authed group — CSRF cookie present, not cacheable.
func TestFormRoutes_KeepCSRFAndAreNotShared(t *testing.T) {
	resp := request(t).setRoute(routenames.Contact).get().assertStatusCode(http.StatusOK)

	var hasCSRF bool
	for _, ck := range resp.Header.Values("Set-Cookie") {
		if strings.Contains(ck, "_csrf") {
			hasCSRF = true
		}
	}
	assert.True(t, hasCSRF, "form pages must keep the CSRF cookie")
	assert.NotContains(t, resp.Header.Get("Cache-Control"), "stale-while-revalidate")
}
