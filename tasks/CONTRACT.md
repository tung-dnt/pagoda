# MIGRATION CONTRACT — ent → sqlc + golang-migrate
Read this fully before editing. Do NOT deviate. Do NOT edit files outside your assigned set.

## What is happening
`entgo.io/ent` is being removed from `github.com/mikestefanello/pagoda` (dir: /Users/tung-dnt/projects/meme-app).
Replaced by **sqlc (postgres/pgx v5)** + **golang-migrate**. The whole `ent/` directory will be deleted.
The DB moves from SQLite to **PostgreSQL**.

## ALREADY DONE (do not recreate; treat as given, read the files)
- `sqlc.yaml` (root)
- `pkg/postgres/migrations/000001_init.{up,down}.sql`
- `pkg/postgres/queries/{users,password_tokens}.sql`
- `pkg/postgres/db/*` — GENERATED, package `pgdb`. NEVER hand-edit.
- `pkg/postgres/migrate.go` — package `postgres`, `Migrate(connection string) error`, `Drop(connection string) error`

## Generated types (package `pgdb`, import `github.com/mikestefanello/pagoda/pkg/postgres/db`)
```go
type User struct { ID int64; Name string; Email string; Password string; Verified bool; Admin bool; CreatedAt time.Time }
type PasswordToken struct { ID int64; Token string; UserID int64; CreatedAt time.Time }

type Querier interface {
	CreatePasswordToken(ctx context.Context, arg CreatePasswordTokenParams) (PasswordToken, error) // {Token string; UserID int64}
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)                            // {Name, Email, Password string; Verified, Admin bool}
	DeletePasswordTokensByUser(ctx context.Context, userID int64) error
	GetPasswordToken(ctx context.Context, arg GetPasswordTokenParams) (PasswordToken, error)       // {ID, UserID int64; CreatedAfter time.Time}
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	SetUserVerified(ctx context.Context, id int64) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)            // {Password string; ID int64}
}
```
`pgdb.New(db DBTX) *Queries`. `GetUserByEmail` and `CreateUser` already `LOWER()` the email in SQL —
callers must NOT call `strings.ToLower` themselves any more.

## Type mapping — apply everywhere
| old | new |
|---|---|
| `*ent.User` | `*pgdb.User` |
| `*ent.PasswordToken` | `*pgdb.PasswordToken` |
| `*ent.Client` / `c.ORM` | `*pgdb.Queries` / `c.Queries` |
| user/token IDs as `int` | **`int64`** (BIGSERIAL) |
| `case *ent.NotFoundError:` | `errors.Is(err, pgx.ErrNoRows)` (`github.com/jackc/pgx/v5`) |
| `case *ent.ConstraintError:` (dup email) | `var pgErr *pgconn.PgError; errors.As(err, &pgErr) && pgErr.Code == "23505"` (`github.com/jackc/pgx/v5/pgconn`) |
| `ent.IsNotFound(err)` | `errors.Is(err, pgx.ErrNoRows)` |

Because sqlc returns *values*, not pointers: `u, err := q.GetUser(...)` gives `pgdb.User`. Take `&u`
when a `*pgdb.User` is needed.

Note the error shape change: ent returned typed errors via `switch err.(type)`. pgx returns sentinel
errors, so those `switch err.(type)` blocks must become `if/else` on `errors.Is`. Rewrite them cleanly —
do not keep an empty `switch`.

## ent hooks are GONE — invariants move to explicit calls
ent had schema hooks that silently (a) lowercased email, (b) bcrypt-hashed `user.password`,
(c) bcrypt-hashed `password_token.token`. Now:
- (a) handled in SQL (`LOWER()`), nothing to do at call sites.
- (b) callers MUST call `authClient.HashPassword(plain)` before `CreateUser` / `UpdateUserPassword`.
- (c) handled inside `AuthClient.GeneratePasswordResetToken`, nothing to do at call sites.

## Post-migration API of pkg/services (being written by the lead — code to THIS signature)
```go
type Container struct {
	Validator     *Validator
	Web           *echo.Echo
	Config        *config.Config
	Cache         *CacheClient
	Database      *pgxpool.Pool // Postgres, application data
	Queries       *pgdb.Queries // sqlc querier
	TasksDatabase *sql.DB       // SQLite -- backlite task queue ONLY
	Files         afero.Fs
	Mail          *MailClient
	Auth          *AuthClient
	Tasks         *backlite.Client
}
// NOTE: Container.ORM no longer exists.

type AuthClient struct{ /* ... */ }
func NewAuthClient(cfg *config.Config, db *pgdb.Queries) *AuthClient
func (c *AuthClient) Login(ctx echo.Context, userID int64) error
func (c *AuthClient) Logout(ctx echo.Context) error
func (c *AuthClient) GetAuthenticatedUserID(ctx echo.Context) (int64, error)
func (c *AuthClient) GetAuthenticatedUser(ctx echo.Context) (*pgdb.User, error)
func (c *AuthClient) CheckPassword(password, hash string) error
func (c *AuthClient) HashPassword(password string) (string, error)   // NEW
func (c *AuthClient) GeneratePasswordResetToken(ctx echo.Context, userID int64) (string, *pgdb.PasswordToken, error)
func (c *AuthClient) GetValidPasswordToken(ctx echo.Context, userID, tokenID int64, token string) (*pgdb.PasswordToken, error)
func (c *AuthClient) DeletePasswordTokens(ctx echo.Context, userID int64) error
func (c *AuthClient) RandomToken(length int) (string, error)
func (c *AuthClient) GenerateEmailVerificationToken(email string) (string, error)
func (c *AuthClient) ValidateEmailVerificationToken(token string) (string, error)
```

## The admin panel
The ent-codegen'd **entity CRUD** admin (`/admin/entity/*`) is being DELETED (user's decision).
The **backlite task dashboard** (`/admin/tasks`) STAYS — it is unrelated to ent. It now uses
`c.TasksDatabase` (the SQLite handle) instead of `c.Database`.

## Ground rules
- `go build ./...` will NOT pass until every agent is done. Don't panic at unresolved symbols in files
  outside your set. DO make sure YOUR files are internally consistent and correctly typed.
- Do not touch `go.mod` / `go.sum` (the lead owns dependency management).
- Do not delete the `ent/` directory (the lead does this last).
- Do not add new third-party dependencies.
- Match surrounding code style: same comment density, same naming, same idiom. Keep every existing
  doc comment that is still accurate; update ones that reference ent/ORM.
- Scope discipline: change what the migration forces. No drive-by refactors.
