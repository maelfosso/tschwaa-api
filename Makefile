POSTGRESQL_URL=postgres://schwaa:123@localhost:5433/schwaa?sslmode=disable

.PHONY: build cover start test test-integration

build:
	docker build -t golangdk-canvas .
cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...

migrate-create:
	migrate create -ext sql -dir storage/migrations -seq $(migration)

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations down