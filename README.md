# Bookmark Manager

Bookmark Manager is a Go backend API project scaffold. In its current state, the implemented application exposes a secure password-generation endpoint through a Gin HTTP server. The repository is structured with API, handler, service, repository, model, and shared package layers, which makes it a starting point for a fuller bookmark-management service.

> Current implementation note: despite the repository name, bookmark CRUD, user accounts, and persistence are not implemented yet. The active feature in the codebase is password generation.

## Project Overview

### What It Does

The application starts an HTTP API server and provides one route:

- `GET /generate-password` returns a randomly generated 12-character password.

The password generator uses Go's `crypto/rand` package and selects characters from a mixed charset containing lowercase letters, uppercase letters, numbers, and special characters.

### Who It Is For

This project is useful for:

- Developers learning or building a layered Go API with Gin.
- Teams starting a backend service and wanting separate API, handler, and service packages.
- API clients that need a simple password-generation endpoint.

### Main Features

- Gin-based HTTP server.
- Environment-based application port configuration.
- Password generation service using cryptographically secure randomness.
- Handler and integration tests using Go's standard testing tools and `testify`.
- Mock-generated service interface for handler tests.
- Reserved folders for future model, repository, and Redis work.

### Why It Is Useful

The codebase demonstrates a small but clean backend architecture:

- `cmd/api` starts the service.
- `internal/api` owns HTTP engine setup and route registration.
- `internal/handler` translates HTTP requests and responses.
- `internal/service` owns business logic.
- Tests validate service behavior, handler responses, and the HTTP endpoint.

## Tech Stack

| Area | Technology |
| --- | --- |
| Language | Go `1.26.2` |
| Backend Framework | Gin `v1.12.0` |
| Configuration | `github.com/kelseyhightower/envconfig` |
| Password Randomness | Go standard library: `crypto/rand`, `math/big` |
| Testing | Go `testing`, `net/http/httptest`, `github.com/stretchr/testify` |
| Mocking | Mockery-generated mock in `internal/service/mocks` |
| Frontend | None currently |
| Database | None currently implemented |
| Shared Packages | Empty Redis package scaffold under `pkg/redis` |
| Build Tooling | Go toolchain; `Makefile` exists but has no targets |

Notes:

- `pkg/redis` currently contains package declarations only. There is no Redis client/config implementation yet.
- `internal/model` and `internal/repository` directories exist but are currently empty.
- `go.mongodb.org/mongo-driver/v2` appears in `go.mod` as an indirect dependency, but the current source code does not use MongoDB.

## Project Structure

```text
bookmark-manager/
|-- cmd/
|   `-- api/
|       `-- main.go
|-- internal/
|   |-- api/
|   |   |-- api.go
|   |   `-- config.go
|   |-- handler/
|   |   |-- genpass.go
|   |   `-- genpass_test.go
|   |-- integration_test/
|   |   `-- genpass_ep_test.go
|   |-- model/
|   |-- repository/
|   `-- service/
|       |-- genpass.go
|       |-- genpass_test.go
|       `-- mocks/
|           `-- genpass.go
|-- pkg/
|   `-- redis/
|       |-- client.go
|       `-- config.go
|-- api.exe
|-- go.mod
|-- go.sum
|-- Makefile
`-- README.md
```

### Important Directories and Files

| Path | Description |
| --- | --- |
| `cmd/api/main.go` | Application entry point. Loads config, creates the API engine, and starts the server. |
| `internal/api/api.go` | Creates the Gin engine, registers routes, starts the HTTP server, and exposes `ServeHTTP` for tests. |
| `internal/api/config.go` | Loads application configuration from environment variables with defaults. |
| `internal/handler/genpass.go` | HTTP handler for `GET /generate-password`. Calls the service and returns JSON. |
| `internal/service/genpass.go` | Password generation business logic and `GenPass` interface. |
| `internal/service/mocks/genpass.go` | Mockery-generated mock for testing handler behavior. |
| `internal/handler/genpass_test.go` | Unit tests for successful and failed password-generation handler responses. |
| `internal/service/genpass_test.go` | Unit tests for generated password lengths. |
| `internal/integration_test/genpass_ep_test.go` | Integration-style HTTP test for the registered API endpoint. |
| `pkg/redis` | Placeholder package for future Redis-related code. |
| `internal/model` | Empty placeholder for future domain models. |
| `internal/repository` | Empty placeholder for future persistence/repository code. |
| `api.exe` | Existing compiled Windows binary artifact. |
| `Makefile` | Present but currently empty. |

## Architecture

The request flow is:

```text
HTTP request
  -> Gin router
  -> handler.GenPass.GeneratePassword
  -> service.GenPass.GeneratePassword
  -> JSON response
```

The handler depends on the `service.GenPass` interface, which allows tests to replace the real service with a mock. This keeps HTTP response tests independent from the random password-generation implementation.

## Features

### Password Generation Endpoint

`GET /generate-password`

The route is registered in `internal/api/api.go`. When called, it invokes `handler.GeneratePassword`, which asks the service to generate a password with a fixed length of `12`.

Successful response:

```json
{
  "password": "generated-value"
}
```

Error response:

```json
{
  "error": "Failed to generate password"
}
```

### Cryptographically Secure Random Passwords

The service in `internal/service/genpass.go` uses `crypto/rand.Int` instead of pseudo-random number generation. Each character is selected from this charset:

```text
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?
```

### Configurable Server Port

The API listens on the configured app port. If no port is provided, it defaults to `8080`.

Because configuration is loaded with `envconfig.Process("api", cfg)`, the preferred variable is:

- `API_APP_PORT`

The code can also fall back to:

- `APP_PORT`

If both are set, `API_APP_PORT` takes precedence.

### Testable HTTP Engine

The API engine implements `ServeHTTP`, allowing endpoint tests to run through `httptest` without starting a real network listener.

## Setup Instructions

### Prerequisites

- Go `1.26.2` or a compatible Go version.
- Network access for the first dependency download.
- Optional: `curl` or another HTTP client for manual endpoint testing.

### Install Dependencies

From the repository root:

```bash
cd bookmark-manager
go mod download
```

### Run the Backend in Development

With the default port:

```bash
go run ./cmd/api
```

With a custom port on macOS/Linux:

```bash
API_APP_PORT=9090 go run ./cmd/api
```

With a custom port on Windows PowerShell:

```powershell
$env:API_APP_PORT = "9090"
go run ./cmd/api
```

The server will listen on:

```text
http://localhost:8080
```

or on the custom port you configured.

### Test the Endpoint

```bash
curl http://localhost:8080/generate-password
```

Example response:

```json
{
  "password": "aB3!example?"
}
```

The returned value will be random and should not match the example exactly.

### Frontend

No frontend application is currently present in this repository.

## API Documentation

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/generate-password` | Generates and returns a 12-character random password. |

### `GET /generate-password`

Generates a password using the service layer.

Request body: none.

Query parameters: none.

Success response:

- Status: `200 OK`
- Body:

```json
{
  "password": "random-12-character-value"
}
```

Failure response:

- Status: `500 Internal Server Error`
- Body:

```json
{
  "error": "Failed to generate password"
}
```

## Environment Variables

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `API_APP_PORT` | No | `8080` | Preferred environment variable for the HTTP server port. |
| `APP_PORT` | No | `8080` | Fallback variable accepted because of the `envconfig:"APP_PORT"` struct tag. |

Only one environment variable is needed to change the port. Prefer `API_APP_PORT`.

## Build and Deployment Guide

### Build a Production Binary

macOS/Linux:

```bash
go build -o api ./cmd/api
./api
```

Windows PowerShell:

```powershell
go build -o api.exe ./cmd/api
.\api.exe
```

### Render

This project can run as a Render Web Service because it is a long-running HTTP server.

Suggested settings:

| Setting | Value |
| --- | --- |
| Runtime | Go |
| Build Command | `go build -o api ./cmd/api` |
| Start Command | `./api` |

Important port note:

- Render recommends binding to the `PORT` environment variable.
- This application currently reads `API_APP_PORT` and `APP_PORT`, not `PORT`.
- Set `API_APP_PORT` to the port Render expects, or update `internal/api/config.go` to support `PORT`.

### Vercel

This repository is not ready for Vercel as-is. The current app is a persistent Gin server started from `cmd/api/main.go`, while Vercel's Go runtime is designed around serverless functions in an `/api` directory with exported HTTP handlers.

To deploy on Vercel, add a Vercel-compatible serverless entry point or adapter. For the current code structure, Render, Fly.io, Railway, a VM, or a container host is a better fit.

### Docker

No `Dockerfile` is currently included. A typical Docker deployment for this app would use a multi-stage Go build and run the compiled binary.

Example Dockerfile:

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.26 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bookmark-manager ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /

COPY --from=build /bookmark-manager /bookmark-manager

ENV API_APP_PORT=8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/bookmark-manager"]
```

Build and run after adding a `Dockerfile`:

```bash
docker build -t bookmark-manager .
docker run --rm -p 8080:8080 -e API_APP_PORT=8080 bookmark-manager
```

## Testing

Run all tests:

```bash
go test ./...
```

Run a specific package:

```bash
go test ./internal/service
go test ./internal/handler
go test ./internal/integration_test
```

Testing currently covers:

- Password generation length behavior.
- Handler success response.
- Handler error response when the service fails.
- End-to-end route registration through the API engine using `httptest`.

Testing libraries and tools:

- Go standard `testing` package.
- Go standard `net/http/httptest`.
- `github.com/stretchr/testify/assert`.
- `github.com/stretchr/testify/mock`.

### Regenerate Mocks

The service interface includes this generator directive:

```go
//go:generate mockery --name GenPass --filename=genpass.go
```

If Mockery is installed, regenerate mocks with:

```bash
go generate ./internal/service
```

## Troubleshooting

### Dependency Download Fails

If `go test ./...` or `go run ./cmd/api` fails during dependency download, confirm network access to the Go module proxy or configure your Go proxy settings.

### Port Does Not Change

Use `API_APP_PORT` first:

```bash
API_APP_PORT=9090 go run ./cmd/api
```

If both `API_APP_PORT` and `APP_PORT` are set, `API_APP_PORT` wins.

### Windows Go Build Cache Issue

If Go cannot initialize its build cache, point the cache to a writable directory:

```powershell
$env:GOCACHE = "$PWD\.gocache"
$env:GOPATH = "$PWD\.gopath"
go test ./...
```

## Contributing

1. Create a focused branch for your change.
2. Keep changes scoped to the relevant package.
3. Add or update tests for behavior changes.
4. Run `go test ./...` before opening a pull request.
5. Run `go mod tidy` after dependency changes.

## License

No license file is currently included in this repository.

## Future Improvements

- Implement bookmark domain models in `internal/model`.
- Add repository implementations in `internal/repository`.
- Complete or remove the Redis scaffold in `pkg/redis`.
- Add bookmark CRUD endpoints if this project is intended to become a full bookmark manager.
- Add a health-check route such as `GET /health`.
- Add an `.env.example` file documenting supported variables.
- Add Makefile targets for common commands like test, run, build, and generate.
- Add a Dockerfile and deployment-specific config once the target platform is chosen.
- Support the common `PORT` environment variable for easier deployment on managed platforms.
- Remove unused dependencies from `go.mod` with `go mod tidy`.
- Add `.gitignore` rules for generated binaries such as `api.exe`.
