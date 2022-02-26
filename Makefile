ENV_FILE := .env
ENV := $(shell cat $(ENV_FILE))

ENV_TEST_FILE := .env.sample
ENV_TEST := $(shell cat $(ENV_TEST_FILE))

APP_NAME := "digicore_v3_backend"

.PHONY: setup
setup:
	go get -v github.com/rubenv/sql-migrate/...
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: gen
gen:
	swag fmt
	swag init

.PHONY: run
run:
	$(ENV) go run main.go

.PHONY: test
test:
	$(ENV_TEST) go test -race -shuffle=on ./...
