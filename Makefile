ifneq (,$(wildcard backend/.env))
    include backend/.env
    export
endif

.PHONY: be fe dev build-be build-fe sqlc migration-new migration-up migration-down seed DB_CHECK

# nextjs commands

fe:
	cd frontend && npm run dev

build-fe:
	cd frontend && npm run build

# go commands

be:
	cd backend && go run cmd/api/main.go

build-be:
	cd backend && go build -o bin/api cmd/api/main.go

# run both backend + frontend
dev:
	npx concurrently \
		"make be" \
		"make fe"

# sqlc commands

sqlc:
	cd backend && sqlc generate

# database commands

DB_CHECK:
ifndef DATABASE_URL
	$(error Error: DATABASE_URL is not set. Make sure your .env file exists and contains it)
endif

# (usage: make migration-new name=add_billing)
migration-new:
ifndef name
	$(error Error: Please provide a migration name. Example: make migration-new name=add_billing_table)
endif
	cd backend && goose -dir sql/migrations create $(name) sql

migration-up: DB_CHECK
	cd backend && goose -dir sql/migrations postgres "$(DATABASE_URL)" up

migration-down: DB_CHECK
	cd backend && goose -dir sql/migrations postgres "$(DATABASE_URL)" down

# (usage: make seed file=mock_admin_audit_logs.sql)
seed: DB_CHECK
ifndef file
	$(error Error: Please provide a seed filename. Example: make seed file=mock_admin_audit_logs.sql)
endif
	@test "$(notdir $(file))" = "$(file)" || \
		(echo "Error: file must be a filename from backend/sql/seeds"; exit 1)
	@test -f "backend/sql/seeds/$(file)" || \
		(echo "Error: seed file not found: backend/sql/seeds/$(file)"; exit 1)
	psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f "backend/sql/seeds/$(file)"
