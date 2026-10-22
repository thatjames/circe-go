---
type: Feature
title: Release builds
description: Pushing a vX.Y.Z tag builds static circe binaries for Linux and macOS and publishes them as a GitHub release.
status: draft
code:
  - .goreleaser.yaml
  - .github/workflows/release.yml
tags: [build, release, github]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Users who do not have Go get a ready binary; versions are tagged so
`go install ...@vX.Y.Z` and the nix-flake package can pin them.

# Behavior
- A pushed tag matching `v*.*.*` starts the release workflow.
- It builds `circe` with `CGO_ENABLED=0` for linux/amd64, linux/arm64,
  darwin/amd64 and darwin/arm64, archives each with `LICENSE` and
  `README.md`, writes a checksums file, and publishes a GitHub release
  for the tag.
- `circe --version` prints the tag the binary was built from (set by
  ldflags); a `go install` build prints the module version from build info.

# Interfaces
- [GitHub Actions CI](/build/ci.md) runs on the same tag.
- [Command surface](/cli/command-surface.md): `--version`.

# Decisions
- **Release builds wanted** (user, 2026-10-22).
- **goreleaser, four targets** (user, 2026-10-22): linux and darwin,
  amd64 and arm64.
