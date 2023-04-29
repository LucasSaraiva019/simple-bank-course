createdb:
	docker exec -it golang-database-1 createdb --username=example --owner=example simple_bank

dropdb:
	docker exec -it golang-database-1 dropdb --username=example simple_bank

migrateup: 
	migrate -path db/migration -database "postgresql://example:example_pass@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://example:example_pass@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker pull kjconroy/sqlc
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

test: 
	go test -v -cover ./...

.PHONY: up createdb dropdb migrateup migratedown sqlc