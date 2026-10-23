---
type: Feature
title: golangci-lint checks
description: A Lint workflow runs golangci-lint v2 with a curated .golangci.yml on every push and pull request and fails on any finding; its badge replaces the retired Go Report Card badge on the README.
status: stable
code:
  - .golangci.yml
  - .github/workflows/lint.yml
tags: [build, ci, github, lint]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T23:00:00Z }
---

# Intent
goreportcard.com has been retired, so its badge on the README is dead
and nothing checks circe's code quality beyond `go vet`. golangci-lint
runs the checks Go Report Card ran (and more) in CI, where a finding
blocks the change instead of only lowering a grade.

# Behavior
- The Lint workflow (`.github/workflows/lint.yml`) runs on every push,
  every pull request and by hand (`workflow_dispatch`), on
  `ubuntu-latest` only, with `CGO_ENABLED=0` and the Go version from
  `go.mod`.
- It runs golangci-lint v2 through `golangci/golangci-lint-action`,
  with the golangci-lint version pinned in the workflow.
- `.golangci.yml` (config `version: "2"`) enables the default linters
  (`errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`) plus
  `misspell`, `gosec` and `revive`, and the formatters `gofmt` and
  `goimports`.
- Any linter finding or unformatted file fails the run; there is no
  baseline of tolerated findings (`new-from-rev` is not set).
- Test files are linted too. Every exclusion carries its reason: a
  rule that is wrong for the whole repository is excluded in
  `.golangci.yml` with a comment; a rule wrong for one line gets
  `//nolint:<linter> // <reason>` on that line.
- Excluded in `.golangci.yml`: the `std-error-handling` preset
  (unchecked `Close`, `Flush` and print errors); gosec G304 (reading
  user-named files is a CLI's job); gosec G104 (errcheck covers it);
  gosec G101 in `_test.go` files (synthetic secrets only).
- A finding whose fix would change behaviour is not fixed in a lint
  change: it gets a `//nolint` naming the backlog plan that fixes it.
- The workflow has `contents: read` only and needs no secret.
- The README shows the Lint workflow's status badge for `main` where the
  Go Report Card badge was; the Go Report Card badge is gone.
- Contributors run the same check locally with `golangci-lint run`
  ([Contributing](/project/contributing.md)).

# Interfaces
- [GitHub Actions CI](/build/ci.md): build, vet and test, the other
  checks a pull request must pass.
- [Dependency updates and vulnerability checks](/build/dependencies.md):
  Dependabot bumps `golangci/golangci-lint-action`, not the pinned
  golangci-lint version.

# Decisions
- **Lint in CI with a badge** (user, 2026-10-24): replaces Go Report
  Card rather than only dropping its badge.
- **Its own workflow on Linux only** (user, 2026-10-24): lint results do
  not depend on the platform, and a separate workflow gets its own
  badge, like Vulnerabilities.
- **A curated `.golangci.yml`** (user, 2026-10-24): the defaults plus
  formatters and a few extra linters, roughly what Go Report Card
  checked.
- **Pinned golangci-lint, bumped by hand** (user, 2026-10-24): a new
  release with new checks cannot break CI unannounced. The release
  binary must be built with a Go of at least the minor version in
  `go.mod`'s `go` line (v2.15.0 is built with go1.27.0, which lints
  `go 1.27.3`), or it refuses to run; check that when raising either.
- **No `make lint` target** (user, 2026-10-24): CI and contributors run
  `golangci-lint run` directly.
