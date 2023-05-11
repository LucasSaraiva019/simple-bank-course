Create docker image:
  docker-compose up --build .

Run:
  docker run --name simplebank --network bank-network -p 8080:8080 -e DB_SOURCE="postgresql://example:example_pass@golang-database-1:5432/simple_bank?sslmode=disable" -e GIN_MODE=release simplebank:latest
