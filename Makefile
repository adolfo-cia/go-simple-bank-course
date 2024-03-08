network:
	docker network create simplebank-network
# connect to existing container to network:
# docker network connect simplebank-network postgres14

postgres-pull:
	docker pull postgres:14.11-alpine3.19

postgres-run:
	docker run --name postgres14 --network=simplebank-network -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:14.11-alpine3.19

postgres-start:
	docker start postgres14

postgres-stop:
	docker stop postgres14

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

# new migration: migrate create -ext sql -dir db/migration -seq init_schema
# add new migration: migrate create -ext sql -dir db/migration -seq add_users
migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -count=1 -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/adolfo-cia/go-simple-bank-course/db/sqlc Store

simplebank-build:
	docker build -t simplebank:latest .

simplebank-run:
	docker run --name simplebank --network=simplebank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@postgres14:5432/simple_bank?sslmode=disable" simplebank:latest
# docker run --name simplebank -p 8080:8080 -e GIN_MODE=release simplebank:latest
# docker run --name simplebank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@172.17.0.3:5432/simple_bank?sslmode=disable" simplebank:latest


.PHONY: network postgres-pull postgres-run postgres-start postgres-stop createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock simplebank-build simplebank-run