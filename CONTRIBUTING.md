# Contributing

## Development setup

```bash
git clone https://github.com/hermes-93/go-utils
cd go-utils
go build ./...   # verify build
make test        # run all tests
make lint        # run golangci-lint (requires golangci-lint installed)
```

## Project layout

```
go-utils/
├── main.go              # CLI entry point — routes os.Args[1] to cmd.*
├── cmd/                 # Flag parsing and I/O wiring per subcommand
│   ├── healthcheck.go
│   ├── portcheck.go
│   ├── logparse.go
│   ├── envcheck.go
│   └── waitfor.go
└── internal/            # Pure logic — no os.Exit, no fmt.Println
    ├── health/
    ├── port/
    ├── logparse/
    └── env/
```

The `internal/` packages have no side effects and are fully unit-testable.
`cmd/` packages are thin: parse flags, call internal, handle exit codes.

## Adding a new subcommand

1. Create `internal/<name>/` with the business logic and a `_test.go` file.
2. Create `cmd/<name>.go` with a `func <Name>(args []string)` entry point.
3. Add the case to the switch in `main.go`.
4. Update `printUsage()` in `main.go`.
5. Document in `README.md`.

## Running tests

```bash
# All tests with race detector
go test ./... -race

# Single package
go test ./internal/health/... -v

# Coverage
make coverage   # opens coverage.html
```

## Code style

- `gofmt` formatting is enforced by CI.
- No external dependencies — standard library only.
- `internal/` packages must not call `os.Exit` or `log.Fatal`.
- Return errors up; let `cmd/` decide how to print and exit.

## Pull requests

1. Fork and create a branch: `git checkout -b feat/my-feature`
2. Make your changes with tests.
3. `make test && make lint` must pass locally.
4. Open a PR against `main`.
