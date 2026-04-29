# go-utils

Go CLI utilities for DevOps automation. A single binary (`devops`) with five focused subcommands for common operational tasks.

## Commands

| Command | Description |
|---------|-------------|
| `healthcheck` | HTTP health check with retry logic |
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

### portcheck

Test TCP port reachability:

```bash
devops portcheck postgres:5432
devops portcheck --timeout 3 redis.internal:6379
```

### logparse

Filter and format JSON structured logs:

```bash
# Show only errors from a log file
devops logparse --level error app.log

# Show info+ from stdin with extra fields
kubectl logs myapp | devops logparse --level info --fields service,trace_id

# Re-emit as JSON
devops logparse --format json --level warn app.log
```

Input lines that are not valid JSON are passed through prefixed with `[raw]`.

### envcheck

Validate required environment variables before starting an application:

```bash
devops envcheck DB_HOST DB_PORT SECRET_KEY
devops envcheck --file .env.example
```

Exit code 1 if any variable is missing or empty.

### waitfor

Block until a service accepts connections (useful in Docker entrypoints):

```bash
# Wait for PostgreSQL (TCP)
devops waitfor --timeout 30 postgres:5432

# Wait for HTTP service to return non-5xx
devops waitfor --http --timeout 60 http://api:8080/health

# In docker-compose
command: sh -c "devops waitfor db:5432 && ./myapp"
```

## Development

```bash
make test       # run tests with race detector
make lint       # golangci-lint
make build      # compile to bin/devops
make docker-build
```

## CI

GitHub Actions runs on every push:
- Build + unit tests on Go 1.22 and 1.23
- golangci-lint
- Docker image build + smoke test

## Architecture

```
go-utils/
├── main.go              # CLI entry point and subcommand routing
├── cmd/                 # Subcommand flag parsing and wiring
│   ├── healthcheck.go
│   ├── portcheck.go
│   ├── logparse.go
│   ├── envcheck.go
│   └── waitfor.go
└── internal/            # Pure business logic, no I/O coupling
    ├── health/          # HTTP checker with retry
    ├── port/            # TCP prober
    ├── logparse/        # JSON log parser
    └── env/             # Env var validator and .env file parser
```

No external dependencies — standard library only.
