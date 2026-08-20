# Contributing to spot

## Prerequisites

- Go 1.24+
- golangci-lint (optional but recommended)

## Quick Start

```sh
git clone https://github.com/heltonmarx/spot.git
cd spot

go build ./...
go test -race ./...
golangci-lint run ./...
```

## Development Workflow

1. Fork the repository on GitHub.
2. Create a feature branch: `git checkout -b my-feature`.
3. Make your changes, keeping them focused and minimal.
4. Format, test, and lint:

   ```sh
   golangci-lint fmt ./...
   go test -race ./...
   golangci-lint run ./...
   ```

5. Open a pull request against `main` describing what changed and why.

## Code Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go).
- Every exported symbol gets a doc comment that explains WHY/when it is used, not just paraphrasing the name.
- Add `ExampleXxx` test functions for exported APIs where useful.
- Use table-driven tests.

See also the pinned `golangci-lint` config in `.golangci.yml`.

## Reporting Issues

When opening an issue at https://github.com/heltonmarx/spot/issues, include:

- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs. actual behavior
