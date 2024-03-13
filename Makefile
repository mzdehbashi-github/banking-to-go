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

exportenvs:
	@echo "Exporting environment variables..."
	export DB_DRIVER=postgres && \
	export DB_SOURCE=postgres://postgres:password@localhost:5432/banking?sslmode=disable && \
	export SERVER_ADDRESS=0.0.0.0:8080 && \
	export PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEUMCfSzDpiDDltDxqtjqVOmtI9eg1\n47Xb5EdsDD0lm4VSKsolb2IyZogLdq5lAsupv4nUeTtE3ixqelkzepDhHg==\n-----END PUBLIC KEY-----" && \
	export PRIVATE_KEY="-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIMFMFGiqgNioubUo4xJAjc3txUTvJbFtH8kb21/+6KUAoAoGCCqGSM49\nAwEHoUQDQgAEUMCfSzDpiDDltDxqtjqVOmtI9eg147Xb5EdsDD0lm4VSKsolb2Iy\nZogLdq5lAsupv4nUeTtE3ixqelkzepDhHg==\n-----END EC PRIVATE KEY-----" && \
	echo "Environment variables exported successfully."

.PHONY: createdb dropdb migrateup migratedown server mock exportenvs