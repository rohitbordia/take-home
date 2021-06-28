postgres:
	docker run --name postgres12  -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	sleep 5
	docker exec -it postgres12 createdb --username=root --owner=root take-home

dropdb:
	docker exec -it postgres12 dropdb take-home

migrateup:
	sleep 3
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/take-home?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/take-home?sslmode=disable" -verbose down


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	DATABASE_URL=postgresql://root:secret@localhost:5432/take-home?sslmode=disable go run main.go


.PHONY: postgres createdb dropdb migrateup migratedown  test server
