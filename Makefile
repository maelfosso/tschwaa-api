POSTGRESQL_URL=postgresql://tschwaa:tschwaa@localhost/tschwaa?sslmode=disable

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

sqlc:
	sqlc generate

migrate-create:
	migrate create -ext sql -dir storage/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations -verbose up $(filter-out $@,$(MAKECMDGOALS))

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations -verbose down $(filter-out $@,$(MAKECMDGOALS))

migrate-fix:
	migrate -database ${POSTGRESQL_URL} -path storage/migrations force $(filter-out $@,$(MAKECMDGOALS))
