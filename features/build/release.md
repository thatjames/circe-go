---
type: Feature
title: Release builds
description: Pushing a vX.Y.Z tag builds static circe binaries for Linux and macOS, publishes them as a GitHub release, and pushes the Docker image to ghcr.io/thatjames/circe.
status: stable
code:
  - .goreleaser.yaml
  - .github/workflows/release.yml
tags: [build, release, github]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:00:22Z }
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
- Archives are `circe_<version>_<os>_<arch>.tar.gz` (version without
  the `v`), each holding `circe`, `LICENSE` and `README.md`; the
  checksums file is `checksums.txt` (SHA-256). The release notes are the
  commit subjects since the previous tag.
- Builds are `-trimpath`, stripped (`-s -w`), with file times set to
  the commit time, so the same tag builds the same binaries.
- `circe --version` prints the tag the binary was built from (set by
  ldflags, `-X github.com/thatjames/circe-go/internal/cli.version=<tag>`);
  a `go install` build prints the module version from build info.
- The workflow runs goreleaser v2 with the repository's `GITHUB_TOKEN`
  (`contents: write`); it needs no other secret.
- A second job, `image`, logs in to `ghcr.io` with `GITHUB_TOKEN`
  (`packages: write`), sets up buildx and runs
  `make docker-push IMAGE=ghcr.io/thatjames/circe`: the
  [Docker image](/build/docker-image.md) for linux/amd64 and
  linux/arm64, tagged `<tag>` and `latest`. The two jobs run in
  parallel; one failing does not undo the other.
- `goreleaser check` validates the config; `goreleaser release --snapshot
  --clean` builds the four archives locally into `dist/` (git-ignored)
  without publishing.

# Interfaces
- [GitHub Actions CI](/build/ci.md) runs on the same tag.
- [Command surface](/cli/command-surface.md): `--version`.

# Decisions
- **Release builds wanted** (user, 2026-10-22).
- **goreleaser, four targets** (user, 2026-10-22): linux and darwin,
  amd64 and arm64.
- **Image on ghcr, only for a tag** (user, 2026-10-23). The package is
  named `circe` after the binary; the image's
  `org.opencontainers.image.source` label links it to the repository.
  ghcr creates a new package private; making it public is a one-time
  step in the GitHub package settings.
