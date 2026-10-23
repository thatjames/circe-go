---
type: Feature
title: Make targets
description: A Makefile builds the circe binary (default), the Docker image tagged with the latest git tag (docker, docker-push) and starts the floci dev container (dev); CI calls it.
status: stable
code:
  - Makefile
tags: [build, make]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:00:35Z }
---

# Intent
One entry point builds every artifact the same way on a laptop and in
CI, so CI does not repeat build commands of its own.

# Behavior
- `VERSION` is the latest tag reachable from HEAD
  (`git describe --tags --abbrev=0`), else `dev` when the clone has no
  tags. It can be overridden: `make VERSION=v0.0.1`.
- `make` (target `build`, the default) builds `./circe` from
  `./cmd/circe` with `CGO_ENABLED=0`, `-trimpath` and the same
  ldflags as the [release builds](/build/release.md):
  `-s -w -X github.com/thatjames/circe-go/internal/cli.version=$(VERSION)`.
  `circe --version` prints `VERSION`. `./circe` is git-ignored.
- `make docker` builds the [Docker image](/build/docker-image.md)
  `$(IMAGE):$(VERSION)`, `IMAGE` defaulting to `circe`, passing
  `VERSION` as a build argument, and loads it into the local Docker
  image store (`--load`, needed with a docker-container buildx
  builder). It does not push.
- `make docker-push` builds the image for `PLATFORMS` (default
  `linux/amd64,linux/arm64`) with `docker buildx` and pushes
  `$(IMAGE):$(VERSION)` and `$(IMAGE):latest`. It needs a buildx builder
  that can build several platforms and a registry login; the release
  workflow gives it both.
- `make dev` starts the [floci dev environment](/build/dev-floci.md)
  with `docker compose up -d --wait` and returns once floci is healthy;
  `make dev-down` stops it and keeps its data volume.
- `make clean` removes `./circe` and `dist/`.
- Targets that make no file are `.PHONY`; the binary is always rebuilt.

# Interfaces
- [GitHub Actions CI](/build/ci.md) runs `make`; the
  [release workflow](/build/release.md) runs `make docker-push`.
- [Go module and toolchain](/build/go-module.md).

# Decisions
- **Makefile as the build entry point, called by CI** (user, 2026-10-23).
- **Image tag is the latest git tag** (user, 2026-10-23: "the latest tag
  seen in github"). Tags reach a clone from GitHub, so CI checks out with
  full history and tags (`fetch-depth: 0`); `dev` when there are none.
- goreleaser keeps building release archives; the Makefile mirrors its
  flags rather than calling it, so `make` needs only Go.
