# gox

[![Go Reference](https://pkg.go.dev/badge/github.com/chinayin/gox.svg)](https://pkg.go.dev/github.com/chinayin/gox)
[![Go Report Card](https://goreportcard.com/badge/github.com/chinayin/gox)](https://goreportcard.com/report/github.com/chinayin/gox)
[![CI](https://github.com/chinayin/gox/workflows/CI/badge.svg)](https://github.com/chinayin/gox/actions)

A collection of useful Go extensions and utilities for internal use.

## Features

- **errors** - Enhanced error handling utilities (Planned)
- **strings** - String manipulation helpers (Planned)
- **slices** - Slice operation utilities (Planned)
- **maps** - Map operation utilities (Planned)
- **sync** - Concurrency utilities and safe data structures (Planned)
- **time** - Time and date utilities (Planned)
- **crypto** - Cryptography and hashing utilities (Planned)
- **net** - Network utilities (HTTP client, IP helpers) (Planned)
- **encoding** - Encoding/decoding utilities (JSON, XML) (Planned)
- **config** - Configuration management utilities
- **validator** - Data validation utilities
- **retry** - Retry mechanism with backoff strategies (Planned)
- **cache** - In-memory caching solutions (Planned)
- **log** - Logging utilities

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

