postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it golang-database-1 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it golang-database-1 dropdb --username=root simple_bank

migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	docker pull kjconroy/sqlc
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

test: 
	go test -v -cover ./...

server:
	go run cmd/simple-bank-server/main.go

mock:
	mockgen --package mockdb --destination db/mock/store.go  github.com/LucasSaraiva019/simple-bank-course/db/sqlc Store
	
.PHONY: up createdb dropdb migrateup migratedown sqlc server mock migrateup1 migratedown1