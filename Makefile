DB_URL=postgresql://foody_owner:diegoelmono@localhost:5432/foody_db?sslmode=disable

init_postgres:
	docker run --name foody_db -p 5432:5432 -e POSTGRES_USER=foody_owner -e POSTGRES_PASSWORD=diegoelmono -d postgres:14-alpine

create_db:
	docker exec -it foody_db createdb --username=foody_owner --owner=foody_owner foody_db

start_db:
	docker start foody_db

migrate_up:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down

migrate_up_1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up 1

migrate_down_1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

dev:
	go run ./cmd/server/main.go

.PHONY: init_postgres create_db start_db migrate_up migrate_down migrate_up_1 migrate_down_1 sqlc dev new_migration
