DB_HOST ?= 127.0.0.1
FILE ?= dump.sql # Temporary file to store the database dump

create_postgres:
	docker run --name gin-postgres -p 5431:5432 -e POSTGRES_USER=mole_user -e POSTGRES_PASSWORD=secret -d timescale/timescaledb-ha:pg17
# sleep 3 
# docker exec gin-postgres bash -c "apt update && apt install -y build-essential postgresql-server-dev-16 pgxnclient && pgxn install pg_uuidv7"

remove_postgres:
	docker stop gin-postgres
	docker rm gin-postgres
	
create_db:
	docker exec gin-postgres createdb --username=mole_user --owner=mole_user mole

drop_db:
	docker exec gin-postgres dropdb mole

migrate_up:
	migrate -path db/migration -database "postgresql://mole_user:secret@localhost:5431/mole?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://mole_user:secret@localhost:5431/mole?sslmode=disable" -verbose down

migrate_force:
	migrate -path db/migration -database "postgresql://mole_user:secret@localhost:5431/mole?sslmode=disable" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover $(shell go list ./...)

mock:
	mockgen -source=db/sqlc/querier.go -destination=db/mock/entity.go Entity 

export_db:
	docker exec gin-postgres pg_dump -U mole_user -Fc \
	  --exclude-schema=_timescaledb_internal \
	  --exclude-schema=_timescaledb_catalog \
	  --exclude-schema=_timescaledb_config \
	  mole > $(FILE)

load_db:
	docker exec -i gin-postgres pg_restore -U mole_user  -d mole < $(FILE)

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mock export_db load_db server
