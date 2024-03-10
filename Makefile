createdb:
	docker exec -it database createdb --username=postgres --owner=postgres banking

dropdb:
	docker exec -it database dropdb --if-exists --username=postgres banking

migratenew:
	migrate create -ext sql -dir db/migrations -seq add_users

migrateup:
	migrate -path db/migrations/ -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path db/migrations/ -database "$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate

test:
	go test ./... -v -cover

server:
	go run main.go

mock:
	mockery --dir=db/sqlc --name=Store --output=db/sqlc/mocks --filename=store.go

.PHONY: createdb dropdb migrateup migratedown server mock