# Get the OS name in lowercase (linux, darwin)
OS_SYSNAME := $(shell uname -s | tr A-Z a-z)
# Get the machine architecture (x86_64, arm64)
OS_MACHINE := $(shell uname -m)

# If mac OS, use `macos-arm64` or `macos-x64`
ifeq ($(OS_SYSNAME),darwin)
	OS_SYSNAME = macos
	ifneq ($(OS_MACHINE),arm64)
		OS_MACHINE = x64
	endif
endif

# If Linux, use `linux-x64`
ifeq ($(OS_SYSNAME),linux)
	OS_MACHINE = x64
endif

# The appropriate Tailwind package for your OS will attempt to be automatically determined.
# If this is not working, hard-code the package you want using these options:
# https://github.com/tailwindlabs/tailwindcss/releases/latest
TAILWIND_PACKAGE = tailwindcss-$(OS_SYSNAME)-$(OS_MACHINE)

# Where the golang-migrate migration files live. sqlc reads the same directory as its schema source.
MIGRATIONS_DIR = pkg/postgres/migrations

# Connection string used by the golang-migrate CLI targets. Override to point at another environment,
# ie: make migrate-up DATABASE_URL="postgres://..."
DATABASE_URL ?= postgres://memleak:memleak@localhost:5432/memleak?sslmode=disable

# The golang-migrate CLI needs its database driver selected at compile time via a build tag, which
# `go tool` cannot pass. It is therefore run with `go run -tags pgx5` instead — still pinned to the
# github.com/golang-migrate/migrate/v4 version in go.mod, and using the same pgx/v5 driver the app
# itself registers in pkg/postgres/migrate.go. That driver is registered under the `pgx5://` scheme.
MIGRATE = go run -tags pgx5 github.com/golang-migrate/migrate/v4/cmd/migrate
MIGRATE_URL = $(subst postgres://,pgx5://,$(DATABASE_URL))
ADMIN_EMAIL = admin@home.local

.PHONY: help
help: ## Print make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: deps-install tailwind-install ## Install all dependencies

.PHONY: deps-install
deps-install: ## Download the Go modules and pre-build the tools pinned in go.mod
	go mod download
	go build -o ./tmp/tools/ tool

.PHONY: tailwind-install
tailwind-install: ## Install the Tailwind CSS CLI
	# Download to ./tmp and move into place rather than overwriting the binary in place: on macOS,
	# replacing a running-or-previously-run binary's contents leaves a stale code-signature cache
	# for that inode and the next exec is SIGKILLed. `mv` swaps in a new inode, so it cannot happen.
	mkdir -p tmp
	curl -sLo tmp/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$(TAILWIND_PACKAGE)
	chmod +x tmp/tailwindcss
	mv -f tmp/tailwindcss tailwindcss
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui.js
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui-theme.js

.PHONY: db-up
db-up: ## Start PostgreSQL via Docker
	docker compose up -d postgres

.PHONY: db-down
db-down: ## Stop PostgreSQL
	docker compose down

.PHONY: sqlc-gen
sqlc-gen: ## Generate type-safe Go code from the SQL queries
	go tool sqlc generate

.PHONY: sqlc-vet
sqlc-vet: ## Lint the SQL queries
	go tool sqlc vet

.PHONY: migrate-new
migrate-new: ## Create a new migration (ie, make migrate-new name=add_posts)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATE_URL)" up

.PHONY: migrate-down
migrate-down: ## Roll back the most recent migration
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATE_URL)" down 1

.PHONY: migrate-force
migrate-force: ## Clear a dirty migration state (ie, make migrate-force version=1)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(MIGRATE_URL)" force $(version)

.PHONY: admin
admin: ## Create a new admin user (ie, make admin ADMIN_EMAIL=me@web.com [DATABASE_URL=postgres://...])
	GO_DATABASE_CONNECTION="$(DATABASE_URL)" go run cmd/admin/main.go --email=$(ADMIN_EMAIL)

.PHONY: run
run: ## Run the application
	clear
	go run cmd/web/main.go

.PHONY: watch
watch: ## Run the application and watch for changes with air to automatically rebuild
	clear
	go tool air

.PHONY: test
test: ## Run all tests (requires PostgreSQL, see `make db-up`)
	go test ./...

.PHONY: check-updates
check-updates: ## Check for direct dependency updates
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["

.PHONY: css
css: ## Build and minify Tailwind CSS
	./tailwindcss -i tailwind.css -o public/static/main.css -m

.PHONY: build
build: css ## Build CSS and compile the application binary
	go build -o ./tmp/main ./cmd/web
