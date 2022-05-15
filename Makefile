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
	# First dev on Windows
	# docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate
	# If using Unix-based
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store_data.go github.com/gunhachi/poke-blackmarket/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
