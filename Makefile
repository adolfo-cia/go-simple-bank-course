network:
	docker network create simplebank-network
# connect to existing container to network:
# docker network connect simplebank-network simplebank-db

postgres-pull:
	docker pull postgres:14.11-alpine3.19

postgres-run:
	docker run --name simplebank-db --network=simplebank-network -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:14.11-alpine3.19

postgres-start:
	docker start simplebank-db

postgres-stop:
	docker stop simplebank-db

createdb:
	docker exec -it simplebank-db createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank-db dropdb simple_bank

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
	docker build -t simplebank-api:latest .

simplebank-run:
	docker run --name simplebank-api --network=simplebank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@simplebank-db:5432/simple_bank?sslmode=disable" simplebank-api:latest
# docker run --name simplebank-api -p 8080:8080 -e GIN_MODE=release simplebank:latest
# docker run --name simplebank-api -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@172.17.0.3:5432/simple_bank?sslmode=disable" simplebank:latest
# docker run --name simplebank-api -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" simplebank:latest


.PHONY: network postgres-pull postgres-run postgres-start postgres-stop createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock simplebank-build simplebank-run