---
type: Feature
title: Command surface
description: The circe binary is a cobra CLI with one subcommand per tinkey operation, shared keyset input/output flags, and errors on stderr with exit 1.
status: draft
code:
  - cmd/circe/**
  - internal/cli/root.go
  - internal/cli/root_test.go
tags: [cli, cobra]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:10:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
circe gives every operation of Tink's tinkey tool on the command line,
written in Go on tink-go. Compatibility is with the Tink API and keyset
formats, not with the Java tinkey CLI's command or flag spelling.

# Behavior
- The binary is `circe`, installed with
  `go install github.com/thatjames/circe-go/cmd/circe@latest`.
- `circe` with no subcommand, or `circe help`, prints the command list.
  `circe --version` prints the version ([release builds](/build/release.md)).
- Subcommands: see [create keyset](/keyset/create-keyset.md),
  [key management](/keyset/key-management.md),
  [public keyset](/keyset/public-keyset.md),
  [inspect and convert](/keyset/inspect-convert.md).
- Commands that read a keyset take `--in <path>` (default: stdin) and
  `--in-format json|binary` (default `json`). Commands that write one take
  `--out <path>` (default: stdout) and `--out-format json|binary`
  (default `json`). Reading and writing go through
  [keyset I/O](/keyset/keyset-io.md).
- `--out` refuses to overwrite an existing file and exits 1, so a keyset
  is never lost by a typo.
- Encrypted keysets take `--master-key-uri` and `--credential`
  ([encrypted keysets](/kms/encrypted-keysets.md)).
- Errors go to stderr as one line, `circe: <message>`, exit code 1.
  Usage errors (unknown flag, missing argument) print the usage and exit 2.
- Output on stdout is only the requested data (keyset, listing), so it
  can be piped.

# Interfaces
- `cmd/circe/main.go` calls `cli.Execute()`; `internal/cli` holds the
  cobra root and one file per subcommand.

# Decisions
- **Binary name `circe`** (user, 2026-10-22).
- **cobra** for the CLI (user, 2026-10-22).
- **Free CLI surface** (user, 2026-10-22): command and flag names may
  differ from Java tinkey. Where there is no reason to differ, tinkey's
  names are kept so its docs carry over.
- **Tests per command** (user, 2026-10-22): every subcommand has tests
  that run it with synthetic keysets.
- **`--out` never overwrites; exit 1 for errors, 2 for usage errors**
  (user, 2026-10-22).
