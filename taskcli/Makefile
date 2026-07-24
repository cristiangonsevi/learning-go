include .env

export

MIGRATION_PATH ?= "./db/migrations"

DB_URL="sqlite3://$(DB_PATH)"

check-migrate:
	@which migrate >/dev/null 2>&1 || { \
		echo "Package migrate not installed"; \
		exit 1; \
	}

migrate-up: check-migrate
	@echo "Running migrations for $(DB_PATH)"
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down: check-migrate
	@echo "Rollback migrations"
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) down 1

create-db:
	@test -f $(DB_PATH) || { \
		echo "Creating database"; \
		sqlite3 $(DB_PATH) "VACUUM;"; \
		echo "Database createde"; \
	}
	@test -f $(DB_PAT) && echo "Database already exist"
