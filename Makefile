FILE = .env
ifneq ($(wildcard $(FILE)),)
include .env
export $(shell sed 's/=.*//' .env)
endif

DB = postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}

.PHONY: migration
migration:
	migrate create -ext sql -dir migrations/ $(name)

.PHONY: migrate
migrate:
	migrate -path migrations/ -database "$(DB)?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations/ -database "$(DB)?sslmode=disable" -verbose down
