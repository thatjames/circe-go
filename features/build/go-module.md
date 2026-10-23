---
type: Feature
title: Go module and toolchain
description: circe-go is the pure-Go module github.com/thatjames/circe-go on Go 1.27, depending on tink-go and cobra only through go.mod.
status: stable
code:
  - go.mod
  - doc.go
tags: [build, go]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:57:50Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
One Go module holds the circe CLI and its packages, so it installs with
`go install` and packages in Nix without extra tooling.

# Behavior
- The module path is `github.com/thatjames/circe-go`.
- `go.mod` declares `go 1.27.3`; the minimum supported Go is 1.27.
  The patch level is raised when a Go release fixes a standard library
  vulnerability circe calls, and the Dockerfile's `golang` build image
  moves with it.
- The module builds with `CGO_ENABLED=0`: no package in it, and no
  dependency it pulls in, needs cgo.
- `go.sum` is committed beside `go.mod` once there are dependencies.
- The module root is package `circe` (`doc.go`), holding only the module's
  package comment; the code lives under `cmd/` and `internal/`.

# Interfaces
- Dependencies: `github.com/tink-crypto/tink-go/v2` for all key
  operations, `github.com/spf13/cobra` (with `spf13/pflag`) for the CLI,
  `google.golang.org/protobuf` for the template protos, and for
  [encrypted keysets](/kms/encrypted-keysets.md) the tink-go KMS
  extensions `tink-go-gcpkms/v2`, `tink-go-awskms/v2` and
  `tink-go-hcvault/v2`, with `google.golang.org/api` (credential options)
  and `hashicorp/vault/api` (the Vault client). All build with cgo off.
- Entry point: `cmd/circe`, see [command surface](/cli/command-surface.md).

# Decisions
- **Pure Go, no cgo** (user, 2026-10-22): static binaries, simple
  `go install` and Nix packaging.
- **Go 1.27 minimum** (user, 2026-10-23, replacing the Go 1.26 minimum of
  2026-10-22): keep `go.mod` at `go 1.27.x`. Taken when Dependabot
  proposed a golang 1.27 build image, so the image and `go.mod`
  stay on one Go version.
- **`go 1.27.3`** (2026-10-23): 1.27.2 still has GO-2026-6809 in
  crypto/x509, which govulncheck finds circe calling.
- **`go 1.26.10`** (2026-10-23): fixes crypto/x509 and net/http
  vulnerabilities govulncheck found circe calling in 1.26.8
  ([Dependency updates](/build/dependencies.md)).
- **Platforms: Linux and macOS** (user, 2026-10-22). Windows is not a
  target.
- **MIT licence.**
- **Root package `doc.go`** (2026-10-22): `go vet ./...` and
  `go test ./...` fail on a module with no packages, so the root carries
  a doc-only package and [CI](/build/ci.md) is green before any code lands.
