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
	PGCONNSTRING = "host=db port=5432 user=postgres password=postgres dbname=calendar sslmode=disable"
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

create-db: os
	docker-compose -f docker-compose.yml up -d
delete-db:
	docker-compose -f docker-compose.yml down -v
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
	easyjson -all cmd/calendar/api/httpapi/api.go
generate:
	go generate ./...
up: prepare_gen generate migrate_db_up gen_json build
	calendar.exe --debug run
.PHONY: build, all, fmt, lint, test, run, os,  generate
