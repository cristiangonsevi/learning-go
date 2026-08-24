-include .env
export

MIGRATION_DIR ?= "file://./db/migrations"
DB_URL ?= "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

build:
	go build -o ./tmp/url-shortener ./cmd/url-shortener/main.go

run:
	go run ./cmd/url-shortener/main.go

dev:
	air

check-migrate:
	@which migrate > /dev/null 2>&1 || { \
	echo "Package migrate required"; \
	}

migrate-create: check-migrate
	@read -p "Migration name: " NAME \
	echo "Migration name $(NAME)"; \
	migrate create -ext sql -dir $(MIGRATION_DIR)  -seq $$NAME

migrate-up: check-migrate
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	-source $(MIGRATION_DIR) \
	-database $(DB_URL) \
	up

migrate-down: check-migrate
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	-source $(MIGRATION_DIR) \
	-database $(DB_URL) \
	down 1
