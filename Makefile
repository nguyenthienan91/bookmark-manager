run:
	go run ./cmd/api/main.go

build:
	CGO_ENABLED=0 go build -o bin/api ./cmd/api

swagger:
	swag init -g ./cmd/api/main.go -o ./docs

COVERAGE_EXCLUDE := mock|test|vendor|docs|Makefile|main.go

test:
	go test ./... --coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker build -t bookmark-manager .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env bookmark-manager

dev-run: swagger run

.PHONY: run build swagger test docker-build docker-run dev-run