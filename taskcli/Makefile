include .env

export

MIGRATION_PATH ?= "./db/migrations"

DB_URL?="sqlite3://$(DB_PATH)"

check-migrate:
	@which migrate >/dev/null 2>&1 || { \
		echo "Package migrate not installed"; \
		exit 1; \
	}

migrate-up: check-migrate
	@echo "Running migrations for $(DB_PATH)"
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down: check-migrate
	@echo "Rollback migrations with config $(MIGRATION_PATH) $(DB_URL)"
	migrate -source $(MIGRATION_PATH) -database $(DB_URL) down 1

migrate-create: check-migrate
	@read -p "Migration name: ", NAME; \
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $$NAME

create-db:
	@which sqlite3 >/dev/null 2>&1 || { \
	echo "Sqlite3 es requerido para correr este comando"; \
	}
	@test -f $(DB_PATH) || { \
		echo "Creating database"; \
		sqlite3 $(DB_PATH) "VACUUM;"; \
		echo "Database created"; \
	}
	@test -f $(DB_PAT) && echo "Database already exist"
