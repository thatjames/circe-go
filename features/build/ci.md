---
type: Feature
title: GitHub Actions CI
description: On every push and pull request, GitHub Actions runs go vet and go test with cgo off on Linux and macOS.
status: draft
code:
  - .github/workflows/ci.yml
tags: [build, ci, github]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
The repository is public open source; every change on GitHub is checked
on both supported platforms before it lands.

# Behavior
- Triggers on push to any branch and on pull requests.
- A matrix of `ubuntu-latest` and `macos-latest` runs, with
  `CGO_ENABLED=0` and the Go version from `go.mod`:
  `go build ./...`, `go vet ./...`, `go test ./...`.
- The run fails when any step fails; a `go.mod`/`go.sum` that is not tidy
  (`go mod tidy` changes it) also fails the run.
- Tests use synthetic keysets only, created in the test; no KMS or cloud
  credentials are needed or configured in CI.

# Interfaces
- [Go module and toolchain](/build/go-module.md).
- Tests of each command: see the features in `keyset/` and `kms/`.

# Decisions
- **GitHub Actions** (user, 2026-10-22): the remote is GitHub.
- KMS backends are tested against fakes or tink-go's test helpers, not
  live services: CI holds no secrets.
