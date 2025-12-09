# gox

[![Go Reference](https://pkg.go.dev/badge/github.com/chinayin/gox.svg)](https://pkg.go.dev/github.com/chinayin/gox)
[![Go Report Card](https://goreportcard.com/badge/github.com/chinayin/gox)](https://goreportcard.com/report/github.com/chinayin/gox)
[![Lint](https://github.com/chinayin/gox/workflows/Lint/badge.svg)](https://github.com/chinayin/gox/actions)
[![Test](https://github.com/chinayin/gox/workflows/Test/badge.svg)](https://github.com/chinayin/gox/actions)
[![Build](https://github.com/chinayin/gox/workflows/Build/badge.svg)](https://github.com/chinayin/gox/actions)
[![codecov](https://codecov.io/gh/chinayin/gox/branch/main/graph/badge.svg)](https://codecov.io/gh/chinayin/gox)
[![License](https://img.shields.io/github/license/chinayin/gox)](LICENSE)

A collection of useful Go extensions and utilities for internal use.

## Packages

### Available

- **[cli](cli/)** - Command-line output formatting with adapter support
  - Professional startup banners
  - Cobra adapter for auto-extracting CLI parameters
  - Extensible for other CLI frameworks
- **[config](config/)** - Configuration management utilities
- **[log](log/)** - Logging utilities
- **[validator](validator/)** - Data validation utilities

### Planned

- **errors** - Enhanced error handling utilities
- **strings** - String manipulation helpers
- **slices** - Slice operation utilities
- **maps** - Map operation utilities
- **sync** - Concurrency utilities and safe data structures
- **time** - Time and date utilities
- **crypto** - Cryptography and hashing utilities
- **net** - Network utilities (HTTP client, IP helpers)
- **encoding** - Encoding/decoding utilities (JSON, XML)
- **retry** - Retry mechanism with backoff strategies
- **cache** - In-memory caching solutions

## Installation

```bash
go get github.com/chinayin/gox
```

## Requirements

- Go 1.25 or higher

## Usage

```go
import "github.com/chinayin/gox"
```



## Development

### Code Quality

This project uses [golangci-lint](https://golangci-lint.run/) for code quality checks.

```bash
# Run linters
make lint

# Auto-fix issues
make lint-fix

# Run tests
make test

# Run all checks
make check
```



## Documentation

Full documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

