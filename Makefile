createdb:
	docker exec -it database createdb --username=postgres --owner=postgres banking

dropdb:
	docker exec -it database dropdb --if-exists --username=postgres banking

migrateup:
	migrate -path db/migrations/ -database "postgres://postgres:password@localhost:5432/banking?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations/ -database "postgres://postgres:password@localhost:5432/banking?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test ./... -v -cover

.PHONY: createdb dropdb migrateup migratedown