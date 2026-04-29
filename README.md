# go-utils

[![CI](https://github.com/hermes-93/go-utils/actions/workflows/ci.yml/badge.svg)](https://github.com/hermes-93/go-utils/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hermes-93/go-utils)](https://goreportcard.com/report/github.com/hermes-93/go-utils)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go CLI utilities for DevOps automation. A single binary (`devops`) with five focused subcommands for common operational tasks. No external dependencies — standard library only.

## Commands

| Command | Description |
|---------|-------------|
| `healthcheck` | HTTP health check with configurable retries |
| `portcheck` | TCP port availability probe |
| `logparse` | JSON structured log filter and formatter |
| `envcheck` | Required environment variable validator |
| `waitfor` | Service readiness poller (TCP or HTTP) |

## Install

```bash
go install github.com/hermes-93/go-utils@latest
```

Or build from source:

```bash
git clone https://github.com/hermes-93/go-utils
cd go-utils
make build
./bin/devops --help
```

## Usage

### healthcheck

Check an HTTP endpoint with configurable retries and expected status code:

```bash
devops healthcheck https://api.example.com/health
devops healthcheck --timeout 10 --retries 5 --expected 200 http://localhost:8080/ready
```

Exit code 0 = healthy, 1 = unhealthy.

```
OK  https://api.example.com/health — HTTP 200 — 42ms (attempt 1/1)
```

### portcheck

Test TCP port reachability:

```bash
devops portcheck postgres:5432
devops portcheck --timeout 3 redis.internal:6379
```

```
OPEN   postgres:5432 (7ms)
```

### logparse

Filter and format JSON structured logs:

```bash
# Show only errors from a log file
devops logparse --level error app.log

# Show info+ from stdin with extra fields
kubectl logs myapp | devops logparse --level info --fields service,trace_id

# Re-emit as JSON (passthrough pipeline filter)
devops logparse --format json --level warn app.log | jq .
```

Input lines that are not valid JSON are passed through prefixed with `[raw]`.

```
10:05:32 [ERROR] database connection failed service=auth trace_id=abc123
```

### envcheck

Validate required environment variables before starting an application:

```bash
devops envcheck DB_HOST DB_PORT SECRET_KEY
devops envcheck --file .env.example
```

```
OK      DB_HOST
OK      DB_PORT
MISSING SECRET_KEY
```

Exit code 1 if any variable is missing or empty. See `examples/env.example`.

### waitfor

Block until a service accepts connections — useful in Docker entrypoints:

```bash
# Wait for PostgreSQL (TCP)
devops waitfor --timeout 30 postgres:5432

# Wait for HTTP service to return non-5xx
devops waitfor --http --timeout 60 http://api:8080/health
```

```
Waiting for postgres:5432 (timeout 30s)...
READY postgres:5432
```

See `examples/docker-compose.yml` for a full compose example.

## Architecture

```
go-utils/
├── main.go              # Entry point — routes os.Args[1] to cmd.*
├── cmd/                 # Flag parsing and exit code handling (thin layer)
│   ├── healthcheck.go
│   ├── portcheck.go
│   ├── logparse.go
│   ├── envcheck.go
│   └── waitfor.go
└── internal/            # Pure business logic — no os.Exit, no I/O
    ├── health/          # HTTP checker with retry
    ├── port/            # TCP dialer
    ├── logparse/        # JSON parser with level filter
    └── env/             # Env var checker + .env file parser
```

`internal/` packages are fully unit-testable with no mocking needed.

## Development

```bash
make test       # run tests with race detector + coverage
make lint       # golangci-lint
make build      # compile to bin/devops
make docker-build
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full development guide.

## CI

GitHub Actions runs on every push:

| Job | What it checks |
|-----|---------------|
| Build & Test | Go 1.22 and 1.23, race detector, coverage |
| Lint | golangci-lint with errcheck, gosec, staticcheck |
| Security | gosec standalone scan |
| Docker | Multi-stage build + smoke test |

## License

MIT
