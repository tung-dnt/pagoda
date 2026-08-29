package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/backlite"
	"github.com/spf13/afero"
	"github.com/tung-dnt/meme-app/config"
	"github.com/tung-dnt/meme-app/pkg/log"
	"github.com/tung-dnt/meme-app/pkg/postgres"
	pgdb "github.com/tung-dnt/meme-app/pkg/postgres/db"

	_ "modernc.org/sqlite"
)

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests.
type Container struct {
	// Validator stores a validator
	Validator *Validator

	// Web stores the web framework.
	Web *echo.Echo

	// Config stores the application configuration.
	Config *config.Config

	// Cache contains the cache client.
	Cache *CacheClient

	// Queries stores the sqlc-generated querier for the application database.
	Queries *pgdb.Queries

	// TasksDatabase stores the connection to the SQLite database backing the task queue.
	// Backlite is SQLite-only, so this is deliberately separate from Database.
	TasksDatabase *sql.DB

	// Files stores the file system.
	Files afero.Fs

	// Mail stores an email sending client.
	Mail *MailClient

	// Auth stores an authentication client.
	Auth *AuthClient

	// Tasks stores the task client.
	Tasks *backlite.Client

	// databaseConnection stores the resolved connection string used for Database, which for tests
	// contains the randomly-generated schema this container isolates itself within.
	databaseConnection string
}

// NewContainer creates and initializes a new Container.
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initValidator()
	c.initWeb()
	c.initCache()
	c.initDatabase()
	c.initTasksDatabase()
	c.initFiles()
	c.initAuth()
	c.initMail()
	c.initTasks()
	return c
}

// Shutdown gracefully shuts the Container down and disconnects all connections.
func (c *Container) Shutdown() error {
	// Shutdown the web server.
	webCtx, webCancel := context.WithTimeout(context.Background(), c.Config.HTTP.ShutdownTimeout)
	defer webCancel()
	if err := c.Web.Shutdown(webCtx); err != nil {
		return err
	}

	// Shutdown the task runner.
	taskCtx, taskCancel := context.WithTimeout(context.Background(), c.Config.Tasks.ShutdownTimeout)
	defer taskCancel()
	c.Tasks.Stop(taskCtx)

	// Discard the disposable schema this test run created.
	if c.Config.App.Environment == config.EnvTest {
		if err := postgres.DropSchema(c.databaseConnection); err != nil {
			return err
		}
	}

	// Shutdown the task queue database.
	if err := c.TasksDatabase.Close(); err != nil {
		return err
	}

	// Shutdown the cache.
	c.Cache.Close()

	return nil
}

// initConfig initializes configuration.
func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg

	// Configure logging.
	switch cfg.App.Environment {
	case config.EnvProduction:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

// initValidator initializes the validator.
func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

// initWeb initializes the web framework.
func (c *Container) initWeb() {
	c.Web = echo.New()
	c.Web.HideBanner = true
	c.Web.Validator = c.Validator
}

// initCache initializes the cache.
func (c *Container) initCache() {
	store, err := newInMemoryCache(c.Config.Cache.Capacity)
	if err != nil {
		panic(err)
	}

	c.Cache = NewCacheClient(store)
}

// initDatabase initializes the PostgreSQL database. Migrations are deliberately NOT applied here for
// real environments -- schema changes are rolled out separately via the migrate-* make targets. The
// test environment is the exception: its schema is disposable and created from scratch on each run,
// so it has to be migrated before any queries can run against it.
func (c *Container) initDatabase() {
	c.databaseConnection = c.Config.Database.Connection

	if c.Config.App.Environment == config.EnvTest {
		// Each test binary runs in its own process and `go test` runs packages in parallel, so every
		// container isolates itself within a randomly-named schema rather than sharing one and
		// clobbering the others.
		c.databaseConnection = replaceRand(c.Config.Database.TestConnection)

		if err := postgres.CreateSchema(c.databaseConnection); err != nil {
			panic(fmt.Sprintf("failed to create test schema: %v", err))
		}

		// The schema was just created and is therefore empty. Build it out with the same embedded
		// migrations the migrate-* make targets apply, so tests run against the real schema.
		if err := postgres.Migrate(c.databaseConnection); err != nil {
			panic(fmt.Sprintf("failed to migrate test schema: %v", err))
		}
	}

	pool, err := pgxpool.New(context.Background(), c.databaseConnection)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to ping database: %v", err))
	}

	c.Queries = pgdb.New(pool)
}

// initTasksDatabase initializes the SQLite database used by the task queue.
func (c *Container) initTasksDatabase() {
	var err error
	c.TasksDatabase, err = openSQLite(c.Config.Database.TasksConnection)
	if err != nil {
		panic(fmt.Sprintf("failed to open task queue database: %v", err))
	}
}

// initFiles initializes the file system.
func (c *Container) initFiles() {
	// Use in-memory storage for tests.
	if c.Config.App.Environment == config.EnvTest {
		c.Files = afero.NewMemMapFs()
		return
	}

	fs := afero.NewOsFs()
	if err := fs.MkdirAll(c.Config.Files.Directory, 0755); err != nil {
		panic(err)
	}
	c.Files = afero.NewBasePathFs(fs, c.Config.Files.Directory)
}

// initAuth initializes the authentication client.
func (c *Container) initAuth() {
	c.Auth = NewAuthClient(c.Config, c.Queries)
}

// initMail initialize the mail client.
func (c *Container) initMail() {
	var err error
	c.Mail, err = NewMailClient(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create mail client: %v", err))
	}
}

// initTasks initializes the task client.
func (c *Container) initTasks() {
	var err error
	// Backlite is SQLite-only, so the queue always runs against its own database rather than the
	// PostgreSQL pool used for application data.
	c.Tasks, err = backlite.NewClient(backlite.ClientConfig{
		DB:              c.TasksDatabase,
		Logger:          log.Default(),
		NumWorkers:      c.Config.Tasks.Goroutines,
		ReleaseAfter:    c.Config.Tasks.ReleaseAfter,
		CleanupInterval: c.Config.Tasks.CleanupInterval,
	})

	if err != nil {
		panic(fmt.Sprintf("failed to create task client: %v", err))
	}

	if err = c.Tasks.Install(); err != nil {
		panic(fmt.Sprintf("failed to install task schema: %v", err))
	}
}

// openSQLite opens the SQLite database used by the task queue.
func openSQLite(connection string) (*sql.DB, error) {
	// Helper to automatically create the directories that the specified sqlite file
	// should reside in, if one.
	d := strings.Split(connection, "/")
	if len(d) > 1 {
		dirpath := strings.Join(d[:len(d)-1], "/")

		if err := os.MkdirAll(dirpath, 0755); err != nil {
			return nil, err
		}
	}

	return sql.Open("sqlite", replaceRand(connection))
}

// replaceRand substitutes $RAND in a connection string with a random value, which is used to give
// each test process its own isolated database or schema.
func replaceRand(connection string) string {
	return strings.Replace(connection, "$RAND", fmt.Sprint(rand.Int()), 1)
}
