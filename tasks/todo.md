# ent → sqlc + golang-migrate migration (meme-app / pagoda)

## Decisions (confirmed with user 2026-08-25)
- **Engine: PostgreSQL** (pgx/v5), mirroring `~/projects/gokit/sqlc.yaml`.
- **Drop the ent-generated entity CRUD admin panel** entirely.
- Consequence (flagged): `backlite` (task queue) is SQLite-only (`modernc.org/sqlite`).
  → Container keeps a *second*, small SQLite connection purely for tasks.
  → The `/admin/tasks` backlite dashboard **stays**; only `/admin/entity/*` goes.
- Consequence: tests now need a real Postgres (no more in-memory memdb).

## Pinned contract (do not diverge)
- Layout: `sqlc.yaml` (root), `pkg/postgres/{migrations,queries,db}`, package `pgdb`.
- IDs: `BIGSERIAL` → Go `int64` (was ent `int`). Update all call sites.
- `*ent.User` → `*pgdb.User`; `*ent.PasswordToken` → `*pgdb.PasswordToken`.
- `*ent.NotFoundError` → `errors.Is(err, pgx.ErrNoRows)`.
- `*ent.ConstraintError` (dupe email) → `*pgconn.PgError` with `Code == "23505"`.
- ent hooks removed → invariants move explicitly:
  - email lowercasing enforced in SQL (`LOWER($n)`).
  - bcrypt hashing of password + reset token done in `pkg/services` before insert.
- `Container.ORM *ent.Client` → `Container.Queries *pgdb.Queries` + `Container.Database *pgxpool.Pool`.
- `Container.TasksDatabase *sql.DB` (SQLite) for backlite.

## Tasks
- [x] P1 Foundation (main agent): sqlc.yaml, migrations, queries, `sqlc generate`, embedded migrate runner
- [x] P2a Agent: services layer (container.go, auth.go, auth_test.go, services_test.go)
- [x] P2b Agent: admin removal (handlers/admin.go, ui/pages/admin_entity.go, ui/forms/admin_entity*.go, middleware/entity.go, routenames, layouts/primary.go, context keys)
- [x] P2c Agent: remaining consumers (handlers/auth.go, middleware/auth.go, ui/request.go+test, pkg/tests/tests.go, cmd/admin/main.go)
- [x] P3 Teardown (main): delete `ent/`, Makefile targets, config.yaml, .air.toml, docker-compose for PG, `go mod tidy`
- [x] P4 Verify: `go build ./...`, `go vet ./...`, `go test ./...` against a live Postgres

## Review

Done. `go vet ./...` and the full test suite pass against a live Postgres.
51 test functions before, 51 after -- none added, removed or weakened.

### CORRECTION: parallel-test race (found by the ui-tests agent, confirmed, fixed)
My first pass had `initDatabase` call `postgres.Drop(TestConnection)` on every container. Each test
package is its own process and `go test` runs packages in parallel, so packages dropped the schema
out from under each other. My initial "all tests pass" was timing luck -- it reproduced 3/3 with
`go test ./pkg/services/... ./pkg/middleware/... ./pkg/handlers/...` (FK violation 23503).

Fix: each test process now isolates itself in a randomly-named Postgres schema via `search_path`,
reusing the `$RAND` convention the SQLite test DB already used. `postgres.Drop` was replaced by
`CreateSchema` / `DropSchema`; the schema is created on startup and dropped on shutdown. Verified
5/5 green afterwards, still running packages in parallel (no `-p 1` serialization needed), and
0 leaked `test_*` schemas left behind.

### Verified at runtime (not just compiled)
- Embedded migrations apply and are idempotent; Drop->Migrate cycle works (used by the test env).
- Schema matches intent (BIGSERIAL ids, unique email, FK w/ ON DELETE CASCADE, timestamptz).
- `make admin` creates a user with a real bcrypt hash (`$2a$10$`, 60 chars) -- NOT plaintext.
  This was the highest-risk part of dropping ent's hooks.
- HTTP register with `E2E.User@Example.COM` stored as `e2e.user@example.com` (SQL LOWER() works).
- HTTP login with different casing + correct password -> authenticated session.
- Wrong password -> "Invalid credentials". Duplicate email -> the pgconn 23505 branch fires
  and redirects to login with the expected message.
- `/admin/tasks` 200 for admin, 401 for non-admin and anon; `/admin/entity/user` now 404.
- Sidebar no longer renders the "Entities" section.
- Password-reset tokens stored bcrypt-hashed; round-trip + expiry covered by the ported
  `TestAuthClient_GetValidPasswordToken`.

### Left alone deliberately
- `PAGODA.md` (untracked, 69KB) is the upstream Pagoda README and still documents the Ent data
  layer. Not rewritten -- it is a reference copy the user added. `README.md` now carries the
  authoritative description of the new setup and points this out.

### Module rename + an unexplained file deletion (resolved)
Mid-session the module was renamed to `github.com/tung-dnt/meme-app`. I initially attributed this to
session meme-app-9f; that was WRONG -- it confirmed it never touched this repo. The rename is now
complete (go.mod + all imports) and the build is clean. Author unknown.

The same event deleted `cmd/admin/main.go` from the worktree. Neither this session nor meme-app-9f
deleted it. It was unrecoverable from git (HEAD/index hold only the pre-migration Ent version, which
no longer compiles), so it was restored from this session's copy of the migrated version, then
re-verified: `make admin` runs and writes a real bcrypt hash ($2a$10$, 60 chars).

Audit of every other deletion in `git status`: the whole staged `ent/**` tree is the intentional
removal, and `pkg/ui/{forms/admin_entity.go, forms/admin_entity_delete.go, pages/admin_entity.go}`
are the intentional admin-panel removal. Nothing else was lost.

The rename also broke import ordering in 11 files (the new path sorts differently); `gofmt -w`
applied, tree is gofmt-clean again.

### Follow-ups the user may want
- Session cookies store the user ID; it changed from `int` to `int64`, so sessions issued before
  this change will fail the type assertion in `GetAuthenticatedUserID`. Harmless in dev; worth a
  session-store flush if this were ever deployed with real users.
- `backlite` remains on SQLite (`dbs/tasks.db`). If a single datastore is wanted, the task queue
  would need replacing.
