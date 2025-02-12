create_postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=mole_user -e POSTGRES_PASSWORD=secret -d postgis/postgis:16-3.5

remove_postgres:
	docker stop postgres
	docker rm postgres
	
create_db:
	docker exec postgres createdb --username=mole_user --owner=mole_user mole

drop_db:
	docker exec postgres dropdb mole

migrate_up:
	migrate -path db/migration -database "postgresql://mole_user:secret@127.0.0.1:5432/mole?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://mole_user:secret@127.0.0.1:5432/mole?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover $(shell go list ./...)

FILE ?= dump.sql

export_db:
	docker exec postgres pg_dump -U mole_user mole > $(FILE)

load_db:
	cat $(FILE) | docker exec -i postgres psql -U mole_user mole

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
