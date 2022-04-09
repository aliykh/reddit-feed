MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} cmd/main.go

.PHONY: swag
swag-init:
	swag init -g cmd/main.go -o api/docs



MIGRATE := migrate -path=./migrations -database=mongodb://root:rootpassword@localhost:27017/reddit-feed?authSource=admin


.PHONY: migrate
migrate-up:
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1



