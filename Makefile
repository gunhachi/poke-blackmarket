postgres:
	docker run --name postgredb -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:latest

createdb:
	docker exec -it postgredb createdb --username=root --owner=root pokemart

dropdb:
	docker exec -it postgredb dropdb pokemart

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pokemart?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pokemart?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test