.EXPORT_ALL_VARIABLES:

PGSCHEMAPATH = ../internal/interfaces/repository/db/sql
MIGRATIONPATH = internal/interfaces/repository/db/migrations
BUILDPATH = -o calendar ./cmd/calendar/main.go
COMPOSE_CONVERT_WINDOWS_PATHS=1

os:
ifeq ($(OS),Windows_NT)
PGSCHEMAPATH = ..\internal\interfaces\repository\db\sql
MIGRATIONPATH="internal\interfaces\repository\db\migrations"
BUILDPATH=-o calendar.exe cmd\calendar\main.go
endif

create-db: os
	docker-compose -f configs/docker-compose.yml up -d
delete-db:
	docker-compose -f configs/docker-compose.yml down -v
tidy:
	go mod tidy
run:
	go run cmd/calendar/main.go -d
build:
	go build $(BUILDPATH)
test:
	go test -v -race ./...
lint: fmt tidy
	golangci-lint run ./...
fmt:
	go fmt ./...
migrate_db_up: create-db
	timeout 3
	goose -dir $(MIGRATIONPATH) postgres "host=localhost port=5432 user=igor password=igor dbname=calendar sslmode=disable" up
goose_down:
	goose -dir $(MIGRATIONPATH) postgres "host=localhost port=5432 user=igor password=igor dbname=calendar sslmode=disable" down
migrate_db_down: goose_down	delete-db

all: migrate_db_up build
	calendar.exe -d
.PHONY: build, all, fmt, lint, test, run, os
