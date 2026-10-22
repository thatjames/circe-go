---
type: Feature
title: Go module and toolchain
description: circe-go is the pure-Go module github.com/thatjames/circe-go on Go 1.26, depending on tink-go and cobra only through go.mod.
status: stable
code:
  - go.mod
tags: [build, go]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
One Go module holds the circe CLI and its packages, so it installs with
`go install` and packages in Nix without extra tooling.

# Behavior
- The module path is `github.com/thatjames/circe-go`.
- `go.mod` declares `go 1.26.8`; the minimum supported Go is 1.26.
- The module builds with `CGO_ENABLED=0`: no package in it, and no
  dependency it pulls in, needs cgo.
- `go.sum` is committed beside `go.mod` once there are dependencies.
- Built today: only `go.mod`. No packages, no dependencies yet.

# Interfaces
- Planned dependencies: `github.com/tink-crypto/tink-go/v2` for all key
  operations, `github.com/spf13/cobra` for the CLI, and the tink-go KMS
  extensions for [encrypted keysets](/kms/encrypted-keysets.md).
- Entry point (planned): `cmd/circe`, see
  [command surface](/cli/command-surface.md).

# Decisions
- **Pure Go, no cgo** (user, 2026-10-22): static binaries, simple
  `go install` and Nix packaging.
- **Go 1.26 minimum** (user, 2026-10-22): keep `go.mod` at `go 1.26.x`.
- **Platforms: Linux and macOS** (user, 2026-10-22). Windows is not a
  target.
- **MIT licence.**
