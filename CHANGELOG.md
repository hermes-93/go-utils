# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `healthcheck --json` flag for machine-readable output
- `portcheck` now accepts multiple addresses in a single invocation
- `logparse` now handles logfmt (`key=value`) format in addition to JSON
- `envcheck --quiet` flag to suppress OK lines and only print failures
- MIT LICENSE file
- `.editorconfig` for consistent editor settings
- GitHub issue templates (bug report, feature request)

## [1.0.0] - 2026-04-29

### Added
- `healthcheck` — HTTP health check with configurable retries and timeout
- `portcheck` — TCP port availability probe
- `logparse` — JSON structured log filter with level filtering and field extraction
- `envcheck` — required environment variable validator with `.env.example` support
- `waitfor` — TCP/HTTP service readiness poller for container entrypoints
- Multi-stage Dockerfile producing a `scratch`-based image
- GitHub Actions CI: build+test on Go 1.22/1.23, golangci-lint, gosec, Docker smoke test
- GoReleaser configuration for cross-platform binary distribution (linux/darwin/windows × amd64/arm64)

[Unreleased]: https://github.com/hermes-93/go-utils/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/hermes-93/go-utils/releases/tag/v1.0.0
