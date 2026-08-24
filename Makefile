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
DATABASE_URL ?= postgres://pagoda:pagoda@localhost:5432/pagoda?sslmode=disable

.PHONY: help
help: ## Print make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: sqlc-install migrate-install air-install tailwind-install ## Install all dependencies

.PHONY: tailwind-install
tailwind-install: ## Install the Tailwind CSS CLI
	curl -sLo tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$(TAILWIND_PACKAGE)
	chmod +x tailwindcss
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui.js
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui-theme.js

.PHONY: sqlc-install
sqlc-install: ## Install the sqlc code generator
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: migrate-install
migrate-install: ## Install the golang-migrate CLI
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: air-install
air-install: ## Install air
	go install github.com/air-verse/air@latest

.PHONY: db-up
db-up: ## Start PostgreSQL via Docker
	docker compose up -d postgres

.PHONY: db-down
db-down: ## Stop PostgreSQL
	docker compose down

.PHONY: sqlc-gen
sqlc-gen: ## Generate type-safe Go code from the SQL queries
	sqlc generate

.PHONY: sqlc-vet
sqlc-vet: ## Lint the SQL queries
	sqlc vet

.PHONY: migrate-new
migrate-new: ## Create a new migration (ie, make migrate-new name=add_posts)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Roll back the most recent migration
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

.PHONY: migrate-force
migrate-force: ## Clear a dirty migration state (ie, make migrate-force version=1)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

.PHONY: admin
admin: ## Create a new admin user (ie, make admin email=myemail@web.com)
	go run cmd/admin/main.go --email=$(email)

.PHONY: run
run: ## Run the application
	clear
	go run cmd/web/main.go

.PHONY: watch
watch: ## Run the application and watch for changes with air to automatically rebuild
	clear
	air

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
