ifneq (,$(wildcard backend/.env))
    include backend/.env
    export
endif

.PHONY: be fe dev build-be build-fe sqlc migration-new migration-up migration-down DB_CHECK

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