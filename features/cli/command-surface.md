---
type: Feature
title: Command surface
description: The circe binary is a cobra CLI with one subcommand per tinkey operation, shared keyset input/output flags, and errors on stderr with exit 1.
status: stable
code:
  - cmd/circe/**
  - internal/cli/root.go
  - internal/cli/root_test.go
  - internal/cli/keyset_flags.go
  - internal/cli/keyset_flags_test.go
tags: [cli, cobra]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:45:23Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
circe gives every operation of Tink's tinkey tool on the command line,
written in Go on tink-go. Compatibility is with the Tink API and keyset
formats, not with the Java tinkey CLI's command or flag spelling.

# Behavior
- The binary is `circe`, installed with
  `go install github.com/thatjames/circe-go/cmd/circe@latest`.
- `circe` with no subcommand, or `circe help`, prints the command list
  on stdout and exits 0.
- `circe --version` prints `circe <version>`: the version set at link
  time (`-X github.com/thatjames/circe-go/internal/cli.version=<tag>`,
  [release builds](/build/release.md)), else the module version from the
  build info (`go install ...@vX.Y.Z`; a build in a checkout gets Go's
  VCS pseudo-version), else `devel`.
- No shell-completion command (cobra's default `completion` is off).
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
  is never lost by a typo. The output is built in memory first and
  written only when the command succeeds, so a failing command writes
  nothing, not even an empty file. A file `--out` creates has mode 0600:
  it may hold secret key material.
- Encrypted keysets take `--master-key-uri` and `--credential`
  ([encrypted keysets](/kms/encrypted-keysets.md)).
- Errors go to stderr as one line, `circe: <message>`, exit code 1.
  Usage errors (unknown command or flag, a bad flag value such as
  `--in-format yaml`, a missing required flag, a stray argument) print
  `circe: <message>` and the command's usage on stderr and exit 2.
- Output on stdout is only the requested data (keyset, listing), so it
  can be piped.

# Interfaces
- `cmd/circe/main.go` exits with `cli.Execute()`; `internal/cli` holds the
  cobra root (`root.go`), the shared `--in*`/`--out*` flags
  (`keyset_flags.go`) and one file per subcommand.
- `cli.Run(args, stdin, stdout, stderr) int` runs circe in-process; the
  tests of every command use it.
- `addInFlags` and `addOutFlags` take the command's KMS flags (or nil);
  with a master key URI they read and write encrypted keysets
  ([encrypted keysets](/kms/encrypted-keysets.md)).

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
