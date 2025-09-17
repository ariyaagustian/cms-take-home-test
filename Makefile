SHELL := /bin/bash

-include .env
.EXPORT_ALL_VARIABLES:

.PHONY: dev migrate-up migrate-down seed test e2e build up down fmt lint

dev:
	docker compose up -d db redis minio imagor
	cd server && go run ./cmd/api

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

seed:
	cd server && go run ./cmd/seed

test:
	cd server && go test ./... -cover

e2e:
	cd tests/e2e && go test -v

build:
	docker build -t cms-gin-server:latest ./server

up:
	docker compose up -d

down:
	docker compose down

fmt:
	gofmt -s -w ./server