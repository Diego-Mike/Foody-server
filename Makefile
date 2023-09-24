init_postgres:
	docker run --name foody_db -p 5432:5432 -e POSTGRES_USER=foody_owner -e POSTGRES_PASSWORD=diegoelmono -d postgres:14-alpine

create_db:
	docker exec -it foody_db createdb --username=foody_owner --owner=foody_owner foody_db

start_db:
	docker start foody_db

migrate_up_1:
	migrate -path internal/db/migrations -database "postgresql://foody_owner:diegoelmono@localhost:5432/foody_db?sslmode=disable" -verbose up 1

migrate_down_1:
	migrate -path internal/db/migrations -database "postgresql://foody_owner:diegoelmono@localhost:5432/foody_db?sslmode=disable" -verbose down 1

migrate_up_2:
	migrate -path internal/db/migrations -database "postgresql://foody_owner:diegoelmono@localhost:5432/foody_db?sslmode=disable" -verbose up 2

migrate_down_2:
	migrate -path internal/db/migrations -database "postgresql://foody_owner:diegoelmono@localhost:5432/foody_db?sslmode=disable" -verbose down 2

sqlc:
	sqlc generate

dev:
	go run ./cmd/server/main.go

.PHONY: init_postgres create_db start_db migrate_up_1 migrate_down_1 migrate_up_2 migrate_down_2 sqlc dev
