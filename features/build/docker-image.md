---
type: Feature
title: Docker image
description: A multi-stage Dockerfile builds a static circe binary for linux/amd64 and linux/arm64 and ships it alone on distroless static (non-root, CA certificates), entrypoint circe.
status: stable
code:
  - Dockerfile
  - .dockerignore
tags: [build, docker]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:00:35Z }
---

# Intent
Run circe where Go is not installed, for example in a pipeline step that
needs to create or convert a keyset.

# Behavior
- Build stage: `golang` at the Go version of `go.mod`, run on the
  build host's platform and cross-compiling to the target's
  `TARGETOS`/`TARGETARCH` (no emulation), module download
  cached as its own layer, then `go build` with `CGO_ENABLED=0`,
  `-trimpath` and the release ldflags; the `VERSION` build argument
  (default `dev`) is what `circe --version` prints.
- Final stage: `gcr.io/distroless/static-debian12:nonroot`, holding only
  `/usr/local/bin/circe`; it runs as the non-root user and has CA
  certificates, so the KMS backends can reach their HTTPS endpoints.
  It carries the label
  `org.opencontainers.image.source=https://github.com/thatjames/circe-go`.
- `ENTRYPOINT ["circe"]`: `docker run --rm circe:<tag> create-keyset
  --key-template AES256_GCM` prints a keyset; keysets in and out go
  through stdin/stdout or a mounted directory. The process runs as uid
  65532, so mounted files must be readable (and, for `--out`, the
  directory writable) by it.
- `.dockerignore` keeps `.git`, `dist/`, `./circe`, `features/` and
  editor/agent files out of the build context.
- Built by `make docker` ([Make targets](/build/makefile.md)); pushed
  to `ghcr.io/thatjames/circe` only for a release tag
  ([Release builds](/build/release.md)).

# Interfaces
- [Command surface](/cli/command-surface.md): the image's entrypoint.
- [KMS-encrypted keysets](/kms/encrypted-keysets.md): credentials reach
  the container through environment variables or a mounted file, never
  baked into the image.

# Decisions
- **distroless static, not scratch**: the same size class, with CA
  certificates, `/etc/passwd` and a non-root user already in place.
- **ghcr.io, release tags only** (user, 2026-10-23).
