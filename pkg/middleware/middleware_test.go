package middleware

import (
	"os"
	"testing"

	"github.com/tung-dnt/meme-app/config"
	pgdb "github.com/tung-dnt/meme-app/pkg/postgres/db"
	"github.com/tung-dnt/meme-app/pkg/services"
	"github.com/tung-dnt/meme-app/pkg/tests"
)

var (
	c   *services.Container
	usr *pgdb.User
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvTest)

	// Create a new container
	c = services.NewContainer()

	// Create a user
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
