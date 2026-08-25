package tests

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	pgdb "github.com/tung-dnt/meme-app/pkg/postgres/db"
	"github.com/tung-dnt/meme-app/pkg/session"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// NewContext creates a new Echo context for tests using an HTTP test request and response recorder
func NewContext(e *echo.Echo, url string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// InitSession initializes a session for a given Echo context
func InitSession(ctx echo.Context) {
	session.Store(ctx, sessions.NewCookieStore([]byte("secret")))
}

// ExecuteMiddleware executes a middleware function on a given Echo context
func ExecuteMiddleware(ctx echo.Context, mw echo.MiddlewareFunc) error {
	handler := mw(func(c echo.Context) error {
		return nil
	})
	return handler(ctx)
}

// ExecuteHandler executes a handler with an optional stack of middleware
func ExecuteHandler(ctx echo.Context, handler echo.HandlerFunc, mw ...echo.MiddlewareFunc) error {
	return ExecuteMiddleware(ctx, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			run := handler

			for _, w := range mw {
				run = w(run)
			}

			return run(ctx)
		}
	})
}

// AssertHTTPErrorCode asserts an HTTP status code on a given Echo HTTP error
func AssertHTTPErrorCode(t *testing.T, err error, code int) {
	httpError, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, code, httpError.Code)
}

// CreateUser creates a random user entity
func CreateUser(db *pgdb.Queries) (*pgdb.User, error) {
	seed := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), rand.Intn(1000000))

	// The password must be hashed prior to storage.
	pw, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := db.CreateUser(context.Background(), pgdb.CreateUserParams{
		Name:     fmt.Sprintf("Test User %s", seed),
		Email:    fmt.Sprintf("testuser-%s@localhost.localhost", seed),
		Password: string(pw),
		Verified: false,
		Admin:    false,
	})
	if err != nil {
		return nil, err
	}

	return &u, nil
}
