postgres-pull:
	docker pull postgres:14.11-alpine

postgres-run:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:14.11-alpine

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

.PHONY: postgres-pull postgres-run postgres-start postgres-stop createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock