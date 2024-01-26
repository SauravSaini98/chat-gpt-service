include .env


build:
	go build -o app .

run: build
	./app

# ======== DATABASE COMMAND START ==========

DB_PASSWORD ?= ""
DB_SCHEMA ?= public
PG_CONNECTION := "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) dbname=$(DB_NAME) password=$(DB_PASSWORD) sslmode=disable search_path=$(DB_SCHEMA)"

db_drop:
	dropdb $(PG_CONNECTION) -e

db_create:
	createdb $(PG_CONNECTION) -e

db_migrate_status:
	goose -dir db/migrations postgres $(PG_CONNECTION) status

db_migrate_up:
	goose -dir db/migrations postgres $(PG_CONNECTION) up

db_migrate_down:
	goose -dir db/migrations postgres $(PG_CONNECTION) down

db_migrate: db_migrate_up

.PHONY: new_migration

new_migration:
    goose -dir db/migrations postgres $(PG_CONNECTION) create NAME=$(NAME) sql

# ======== DATABASE COMMAND END ==========

.PHONY: build run db_drop db_create db_migrate_status db_migrate_up db_migrate_down db_migrate new_migration