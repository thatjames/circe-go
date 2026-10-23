---
type: Feature
title: GitHub Actions CI
description: On every push and pull request, GitHub Actions runs go vet, go test and make with cgo off on Linux and macOS.
status: stable
code:
  - .github/workflows/ci.yml
tags: [build, ci, github]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:00:22Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
The repository is public open source; every change on GitHub is checked
on both supported platforms before it lands.

# Behavior
- Triggers on push to any branch and on pull requests.
- A matrix of `ubuntu-latest` and `macos-latest` runs, with
  `CGO_ENABLED=0` and the Go version from `go.mod`:
  `go build ./...` and `make` (the binary, [Make targets](/build/makefile.md)),
  `go vet ./...`, `go test ./...`.
- The checkout fetches full history and tags (`fetch-depth: 0`), so
  `make` stamps the binary with the latest tag.
- CI builds no Docker image: the image is built and pushed only for a
  release tag ([Release builds](/build/release.md)).
- The run fails when any step fails; a `go.mod`/`go.sum` that is not tidy
  (`go mod tidy -diff` reports a change) also fails the run.
- The jobs run on GitHub's `actions/checkout` and `actions/setup-go`,
  with `go-version-file: go.mod`; one platform failing does not cancel
  the other (`fail-fast: false`).
- Tests use synthetic keysets only, created in the test; no KMS or cloud
  credentials are needed or configured in CI.

# Interfaces
- [Go module and toolchain](/build/go-module.md).
- Tests of each command: see the features in `keyset/` and `kms/`.

# Decisions
- **GitHub Actions** (user, 2026-10-22): the remote is GitHub.
- **CI builds through the Makefile** (user, 2026-10-23); `make dev`
  stays local.
- **No image outside a release** (user, 2026-10-23): the container is
  built and pushed only when there is a new version (a tag).
- KMS backends are tested against fakes or tink-go's test helpers, not
  live services: CI holds no secrets.
