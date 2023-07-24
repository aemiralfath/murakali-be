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

.PHONY: create-mock
create-mock:
	mockery --dir=./internal/module/admin --name=UseCase --output=./internal/module/admin/mocks
	mockery --dir=./internal/module/auth --name=UseCase --output=./internal/module/auth/mocks
	mockery --dir=./internal/module/cart --name=UseCase --output=./internal/module/cart/mocks
	mockery --dir=./internal/module/location --name=UseCase --output=./internal/module/location/mocks
	mockery --dir=./internal/module/product --name=UseCase --output=./internal/module/product/mocks
	mockery --dir=./internal/module/seller --name=UseCase --output=./internal/module/seller/mocks
	mockery --dir=./internal/module/user --name=UseCase --output=./internal/module/user/mocks
	mockery --dir=./internal/module/admin --name=Repository --output=./internal/module/admin/mocks
	mockery --dir=./internal/module/auth --name=Repository --output=./internal/module/auth/mocks
	mockery --dir=./internal/module/cart --name=Repository --output=./internal/module/cart/mocks
	mockery --dir=./internal/module/location --name=Repository --output=./internal/module/location/mocks
	mockery --dir=./internal/module/product --name=Repository --output=./internal/module/product/mocks
	mockery --dir=./internal/module/seller --name=Repository --output=./internal/module/seller/mocks
	mockery --dir=./internal/module/user --name=Repository --output=./internal/module/user/mocks

.PHONY: test-coverage
test-coverage:
	go test -failfast -tags=integration -coverprofile=$(COVERAGE_OUTPUT) -covermode=atomic ./internal/module/...
	go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_OUTPUT_HTML)