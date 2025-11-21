# osrs-bot-api

A lightweight Go API for controlling OSRS bots and exposing automation features over HTTP. This repository contains the API server code, Docker configuration, and helpers to run and test the service.

> NOTE: This README was generated from a repository scan. Please review and adapt examples (ports, env names, endpoints) to match your code and intended deployment.

Table of contents
- [Quick overview](#quick-overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Run with Docker (recommended)](#run-with-docker-recommended)
- [Run locally (Go)](#run-locally-go)
- [Running tests and linters](#running-tests-and-linters)
- [API documentation](#api-documentation)
- [Development & contribution](#development--contribution)
- [Security & secrets](#security--secrets)
- [License](#license)
- [Contact](#contact)

Quick overview
-------------
osrs-bot-api is a Go-based HTTP API intended to control, query, or orchestrate Old School RuneScape bot actions. The repo contains the server source in Go and a small Docker setup to containerize the service.

Features
--------
- HTTP API implemented in Go
- Docker-friendly, ready to be containerized
- Intended to be extended with authentication, rate limiting, and metrics

Prerequisites
-------------
- Go 1.20+ (or the Go version you use in this repo)
- Docker & Docker Compose (for containerized runs)
- PostgreSQL (if the application expects a DB; adjust to your DB of choice)

Configuration
-------------
The application reads configuration from environment variables. Example (DO NOT commit secrets):

```env
POSTGRES_USER=admin
POSTGRES_PASSWORD=password
POSTGRES_DB=test_db
PORT=8080
# Add any other environment variables your app requires
```

Important: the repository currently contains a committed `.env` file with sample credentials. Remove secrets and add `.env` to `.gitignore` before pushing or sharing. See the Security section below for more.

Run with Docker (recommended)
-----------------------------
Build and run the container locally:

Build:
```bash
docker build -t osrs-bot-api:local .
```

Run (example):
```bash
docker run --rm -p 8080:8080 \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=supersecret \
  -e POSTGRES_DB=osrs_db \
  osrs-bot-api:local
```

Example docker-compose (create `docker-compose.yml`):
```yaml
version: "3.8"
services:
  api:
    image: osrs-bot-api:local
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    depends_on:
      - db

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

Run locally (Go)
----------------
1. Clone the repo:
```bash
git clone https://github.com/TheGreatBabushka/osrs-bot-api.git
cd osrs-bot-api
```

2. Fetch dependencies and build:
```bash
go mod download
go build ./...
```

3. Set environment variables (or use a local `.env`) and run:
```bash
export POSTGRES_USER=admin
export POSTGRES_PASSWORD=supersecret
export POSTGRES_DB=osrs_db
go run ./cmd/your-app-main
```

Adjust the `go run` target to the actual main package path in the repo (e.g., `./cmd/server`).

Running tests and linters
-------------------------
Run unit tests:
```bash
go test ./...
```

Formatting and linting (recommended):
```bash
gofmt -w .
golangci-lint run
```

Add these checks to CI (GitHub Actions) to block regressions.

API documentation
-----------------
This repository does not currently include a formal OpenAPI/Swagger spec. Consider adding:
- OpenAPI 3 spec (yaml/json) in `/api` or `/docs`
- Swagger UI or Redoc served by the app for interactive docs
- Endpoint examples in this README (once endpoints are known)

As a quick check, run the server and open the root or `/health` endpoint (or whichever endpoints your code exposes) to confirm availability.

Development & contribution
--------------------------
- Please add a CONTRIBUTING.md with PR and branching guidelines before accepting external contributions.
- Use atomic commits with descriptive messages and open feature branches per change.
- Add unit tests for new features and validate them in CI.
- Consider adding GitHub templates: ISSUE_TEMPLATE.md and PULL_REQUEST_TEMPLATE.md.

Security & secrets
------------------
- Remove the committed `.env` file from the repository history if it contains secrets. Steps:
  - Remove the file and add `.env` to `.gitignore`
  - If secrets were pushed to a public repo, rotate credentials immediately.
  - For complete history removal, use tools such as `git filter-repo` or `bfg-repo-cleaner`.
- Never commit production credentials. Use a secrets manager or GitHub Secrets for CI.

Suggested .gitignore additions
------------------------------
```
.env
