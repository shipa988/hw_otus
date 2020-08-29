.EXPORT_ALL_VARIABLES:

COMPOSE_CONVERT_WINDOWS_PATHS=1

PGSCHEMAPATH = ./internal/data/repository/db/sql
PGCONNSTRING = "host=localhost port=5432 user=postgres password=postgres dbname=calendar sslmode=disable"
DSN = host=localhost port=5432 user=igor password=igor dbname=calendar sslmode=disable
MIGRATIONPATH = ./internal/data/repository/db/migrations

BUILD_CALENDAR_PATH = ./cmd/calendar/main.go
BUILD_CALENDAR = calendar
CONFIG_CALENDAR_PATH = cmd/calendar/config/config.yaml

BUILD_SCHEDULER_PATH = ./cmd/calendar_scheduler/main.go
BUILD_SCHEDULER = calendar_scheduler
CONFIG_SCHEDULER_PATH = cmd/calendar_scheduler/config/config.yaml

BUILD_SENDER_PATH = ./cmd/calendar_sender/main.go
BUILD_SENDER = calendar_sender
CONFIG_SENDER_PATH = cmd/calendar_sender/config/config.yaml

BUILDPATH = -o calendar ./cmd/calendar/main.go

RABBIT_QUEUE = topic_calendar
RABBIT_ADDR = amqp://guest:guest@localhost:5672/
RABBITMQ_DEFAULT_USER = guest
RABBITMQ_DEFAULT_PASS = guest

docker_env:
DSN = host=db port=5432 user=igor password=igor dbname=calendar sslmode=disable
RABBIT_ADDR = amqp://guest:guest@rabbitmq:5672/

os:
ifeq ($(OS),Windows_NT)
PGSCHEMAPATH = .\internal\data\repository\db\sql
MIGRATIONPATH=".\internal\data\repository\db\migrations"

BUILD_CALENDAR_PATH =  cmd\calendar\main.go
BUILD_SCHEDULER_PATH = cmd\calendar_scheduler\main.go
BUILD_SENDER_PATH = cmd\calendar_sender\main.go

CONFIG_CALENDAR_PATH = cmd\calendar\configs\config.yaml
CONFIG_SCHEDULER_PATH  = cmd\calendar_scheduler\configs\config.yaml
CONFIG_SENDER_PATH = cmd\calendar_sender\configs\config.yaml

BUILD_CALENDAR = calendar.exe
BUILD_scheduler = calendar_scheduler.exe
BUILD_SENDER = calendar_sender.exe

endif
timeout:
	timeout 5
tidy:
	go mod tidy
run:
	go run cmd/calendar/main.go -d
run_calendar:
	go run $(BUILD_CALENDAR_PATH) --config $(CONFIG_CALENDAR_PATH) --debug run
run_scheduler:
	go run $(BUILD_SCHEDULER_PATH) --config $(CONFIG_SCHEDULER_PATH) --debug run
run_sender:
	go run $(BUILD_SENDER_PATH) --config $(CONFIG_SENDER_PATH) --debug run

build_calendar:
	go build -o $(BUILD_CALENDAR) $(BUILD_CALENDAR_PATH)
build_scheduler:
	go build -o $(BUILD_SCHEDULER) $(BUILD_SCHEDULER_PATH)
build_sender:
	go build -o $(BUILD_SENDER) $(BUILD_SENDER_PATH)
build: build_calendar build_scheduler build_sender

start_calendar:
	$(BUILD_CALENDAR) --config $(CONFIG_CALENDAR_PATH) --debug run
start_scheduler:
	$(BUILD_SCHEDULER) --config $(CONFIG_SCHEDULER_PATH) --debug run
start_sender:
	$(BUILD_SENDER) --config $(CONFIG_SENDER_PATH) --debug run

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
testi:
	go test -p 1 -v ./tests -tags=integration
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
prepare_gen:
	go install github.com/mailru/easyjson
	go get \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
        github.com/golang/protobuf/protoc-gen-go
prepare_migrate:
	go get -u github.com/pressly/goose/cmd/goose@v2.7.0-rc4

lint: prepare_lint fmt tidy
	golangci-lint run ./...
fmt:
	go fmt ./...

docker-up-dev: os
	docker-compose -f docker-compose-dev.yml up -d
docker-down-dev:
	docker-compose -f docker-compose-dev.yml down -v --remove-orphans

docker-up: os
	docker-compose -f docker-compose.yml up -d
docker-down:
	docker-compose -f docker-compose.yml down
docker-down-hard:
	docker-compose -f docker-compose.yml down -v --remove-orphans
	docker rmi hw12_13_14_15_calendar_calendar:latest
	docker rmi hw12_13_14_15_calendar_scheduler:latest
	docker rmi hw12_13_14_15_calendar_sender:latest

migrate_db:timeout goose_up

goose_up:
	goose -dir $(MIGRATIONPATH) postgres $(PGCONNSTRING) up
goose_down:
	goose -dir $(MIGRATIONPATH) postgres $(PGCONNSTRING) down

gen_json:
	easyjson -all cmd/calendar/api/httpapi/api.go
generate:
	go generate ./...

up-dev: prepare_gen generate gen_json docker-up-dev prepare_migrate migrate_db build start_calendar
down-dev: docker-down-dev

up-prod: docker_env prepare_gen generate gen_json docker-up
down-prod: docker-down

up-ci: docker_env gen_json docker-up prepare_migrate migrate_db
down-ci: docker-down-hard

up: docker_env prepare_gen generate gen_json docker-up prepare_migrate migrate_db
down: docker-down-hard

integration-tests: up-prod testi down-prod
integration-tests-ci: up-ci testi down-ci

.PHONY: build, all, fmt, lint, test, run, os,up,down,gen,tidy, generate