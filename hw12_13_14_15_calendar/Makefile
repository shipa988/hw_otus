.EXPORT_ALL_VARIABLES:

PGSCHEMAPATH = ../internal/data/repository/db/sql
PGCONNSTRING = "host=localhost port=5432 user=postgres password=postgres dbname=calendar sslmode=disable"
MIGRATIONPATH = internal/data/repository/db/migrations
BUILDPATH = -o calendar ./cmd/calendar/main.go
COMPOSE_CONVERT_WINDOWS_PATHS=1

os:
ifeq ($(OS),Windows_NT)
PGSCHEMAPATH = ..\internal\data\repository\db\sql
MIGRATIONPATH="internal\data\repository\db\migrations"
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
	go test -race ./...
testv:
	go test -v -race ./...
prepare_lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
prepare_gen:
	go install github.com/mailru/easyjson
	go get \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
        github.com/golang/protobuf/protoc-gen-go
lint: prepare_lint fmt tidy
	golangci-lint run ./...
fmt:
	go fmt ./...
migrate_db_up: create-db
	timeout 5
	goose -dir $(MIGRATIONPATH) postgres $(PGCONNSTRING) up
goose_down:
	goose -dir $(MIGRATIONPATH) postgres $(PGCONNSTRING) down
migrate_db_down:	delete-db
gen_json:
	easyjson -all api/httpapi/api.go
generate:
	go generate ./...
all: prepare_gen generate migrate_db_up gen_json build
	calendar.exe -d
.PHONY: build, all, fmt, lint, test, run, os,  generate
