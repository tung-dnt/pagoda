package services

import (
	"os"
	"testing"

	"github.com/tung-dnt/meme-app/config"
	pgdb "github.com/tung-dnt/meme-app/pkg/postgres/db"
	"github.com/tung-dnt/meme-app/pkg/tests"

	"github.com/labstack/echo/v4"
)

var (
	c   *Container
	ctx echo.Context
	usr *pgdb.User
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvTest)

	// Create a new container
	c = NewContainer()

	// Create a web context
	ctx, _ = tests.NewContext(c.Web, "/")
	tests.InitSession(ctx)

	// Create a test user
	var err error
	if usr, err = tests.CreateUser(c.Queries); err != nil {
		panic(err)
	}

	// Run tests
	exitVal := m.Run()

	// Shutdown the container
	if err = c.Shutdown(); err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}
