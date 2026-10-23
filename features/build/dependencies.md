---
type: Feature
title: Dependency updates and vulnerability checks
description: Dependabot opens weekly grouped update pull requests for Go modules, GitHub Actions and Docker base images; govulncheck fails a Vulnerabilities workflow on push, pull request and every Monday when circe calls a known-vulnerable symbol.
status: stable
code:
  - .github/dependabot.yml
  - .github/workflows/vuln.yml
tags: [build, ci, github, security, dependencies]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:36:52Z }
---

# Intent
circe handles key material, so a known vulnerability in Go, tink-go or
a KMS SDK must be found and patched without someone remembering to
look.

# Behavior
- Dependabot (`.github/dependabot.yml`) checks weekly and opens pull
  requests for:
  - Go modules (`go.mod` at the root), all updates in one grouped pull
    request, titled `fix(deps): ...`, so they appear under Fixes in the
    release notes ([Release builds](/build/release.md));
  - GitHub Actions used in `.github/workflows/`, grouped, titled
    `chore(deps): ...` (left out of the release notes);
  - Docker base images in the `Dockerfile` (`golang`, distroless),
    titled `chore(deps): ...`.
- Dependabot pull requests go through the same CI and squash merge as
  any other ([Contributing](/project/contributing.md)).
- The Vulnerabilities workflow (`vuln.yml`) runs on every push, every
  pull request, every Monday at 06:00 UTC and by hand
  (`workflow_dispatch`). It runs
  `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` on Linux
  with `CGO_ENABLED=0` and the Go version from `go.mod`.
- govulncheck fails the run (exit 3) only when circe's code reaches a
  vulnerable symbol, in the standard library of that Go version or in a
  required module; vulnerabilities in code circe never calls are listed
  but pass.
- The weekly run catches a vulnerability published after the last
  change; its badge is on the README.
- The workflow has `contents: read` only and needs no secret.

# Interfaces
- [GitHub Actions CI](/build/ci.md): the other checks a pull request
  must pass.
- [Go module and toolchain](/build/go-module.md): a standard library
  vulnerability is fixed by raising the `go` line in `go.mod` and the
  `golang` image in the `Dockerfile` together.

# Decisions
- **A vulnerability patcher** (user, 2026-10-23): Dependabot version
  and security updates, with govulncheck as the check.
- **govulncheck, not only Dependabot alerts**: it reports only
  vulnerabilities circe actually calls, and it covers the Go standard
  library, which Dependabot does not update.
- **`govulncheck@latest`**: the vulnerability database and the tool
  move together; pinning would hide new findings.
- **Go module updates are `fix(deps)`**: they change the shipped
  binary, so users see them in the release notes; action and base image
  bumps do not change circe's behaviour and are chores.
- **Dependabot security updates** are a repository setting (Settings,
  Code security), switched on by hand; `dependabot.yml` only configures
  version updates.
