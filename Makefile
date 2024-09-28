DB_URL=postgresql://root:123@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres-16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -d postgres:16-alpine

sqlc:
	sqlc generate

createdb:
	docker exec -it postgres-16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-16 dropdb simple_bank

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover ./...

.PHONY: network postgres sqlc createdb dropdb migrateup migratedown test