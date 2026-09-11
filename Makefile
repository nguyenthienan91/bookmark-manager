run:
	go run ./cmd/api/main.go

swagger:
	swag init -g ./cmd/api/main.go -o ./docs

COVERAGE_EXCLUDE := mock|test|vendor|docs|Makefile|main.go

test:
	go test ./... --coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html

dev-run: swagger run

.PHONY: run swagger test dev-run