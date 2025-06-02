DB_URL=postgres://root:secret@localhost:5432/session-db?sslmode=disable
MIGRATIONS_DIR := db/migrations

define create_migration
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(1)
endef

network:
	docker network create dnd-network

postgres:
	docker run --name session-db --network dnd-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

restartpg:
	docker restart session-db

createdb:
	docker exec -it session-db createdb --username=root --owner=root session-db

dropdb:
	docker exec -it session-db dropdb session-db

create_migration:
	$(call create_migration,$(filter-out $@,$(MAKECMDGOALS)))

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migrateversion:
	migrate -path db/migrations/ -database "$(DB_URL)" force 1

sqlc:
	sqlc generate

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateversion