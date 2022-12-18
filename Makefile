COVERAGE_OUTPUT=coverage.out
COVERAGE_OUTPUT_HTML=coverage.html
CMD_SERVER=./cmd/api/main.go
CMD_SERVER_OUT=./build/server
GO_BUILD_FLAGS=-trimpath -a -tags "osusergo,netgo" -ldflags '-extldflags=-static -w -s -v'

ENV ?= local

.PHONY: mod
mod:
	go mod download
	go mod vendor

.PHONY: golangci-lint-install
golangci-lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

.PHONY: golangci-lint-run
golangci-lint-run:
	golangci-lint run ./...

.PHONY: golangci-lint
golangci-lint:
	golangci-lint-install golangci-lint-run

.PHONY: compile
compile:
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(CMD_SERVER_OUT) $(CMD_SERVER)

.PHONY: test-coverage
test-coverage:
	go test -race -failfast -tags=integration -coverprofile=$(COVERAGE_OUTPUT) -covermode=atomic ./internal/...
	go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_OUTPUT_HTML)

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir sql/migrations -seq $(migration_name)

.PHONY: migrate-up
migrate-up:
	migrate -path sql/migrations -database "postgres://$(user):$(pwd)@localhost:5431/$(db)?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path sql/migrations -database "postgres://$(user):$(pwd)@localhost:5431/$(db)?sslmode=disable" -verbose down

.PHONY: migrate-force
migrate-force:
	migrate -path sql/migrations -database "postgres://$(user):$(pwd)@localhost:5431/$(db)?sslmode=disable" -verbose force $(version)