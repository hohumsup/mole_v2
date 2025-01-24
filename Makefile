create-postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=mole_user -e POSTGRES_PASSWORD=secret -d postgis/postgis:16-3.5

remove-postgres:
	docker stop postgres
	docker rm postgres
	
createdb:
	docker exec -it postgres createdb --username=mole_user --owner=mole_user mole

dropdb:
	docker exec -it postgres dropdb mole

migrateup:
	migrate -path db/migration -database "postgresql://mole_user:secret@localhost:5432/mole?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://mole_user:secret@localhost:5432/mole?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
