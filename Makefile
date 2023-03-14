POSTGRESQL_URL=postgres://tschwaa:tschwaa@localhost/tschwaa?sslmode=disable

.PHONY: build cover start test test-integration

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main cmd/server/*

build:
	docker build -t tschwaa-api .

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...

migrate-create:
	migrate create -ext sql -dir storage/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations up $(filter-out $@,$(MAKECMDGOALS))

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations down $(filter-out $@,$(MAKECMDGOALS))