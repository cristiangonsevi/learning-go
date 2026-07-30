include .env

export

BIN_FILENAME ?= "cli"
MIGRATION_DIR ?= "./db/migrations"
DB_URL ?= "sqlite3://$(DB_PATH)"

build:
	@echo "Building app"
	@mkdir -p ./bin
	@go build -o ./bin/$(BIN_FILENAME) .

run: build
	@./bin/$(BIN_FILENAME)

check-migrate:
	@which migrate >/dev/null 2>&1 || { \
	echo "Package migrate not installed"; \
	}

db-create: check-migrate
	@which sqlite3 > /dev/null 2>&1 || { \
	echo "Sqlite3 required to procced"; \
	}
	@test -f $(DB_PATH) || { \
	echo "Creating database"; \
	sqlite3 $(DB_PATH) "VACUUM;"; \
	echo "Database created"; \
	}
	@test -f $(DB_PATH) && echo "Database already exist"

db-migrate-create: check-migrate
	@read -p "Migration name: " NAME
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $$NAME

db-migrate-up: check-migrate
	@echo "Running migrations"
	@migrate -database $(DB_URL) -path $(MIGRATION_DIR) up

db-migrate-down:
	@echo "Running rollback"
	@migrate -database $(DB_URL) -path $(MIGRATION_DIR) dow 1
