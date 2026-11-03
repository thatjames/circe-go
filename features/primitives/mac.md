---
type: Feature
title: MAC compute and verify
description: circe mac compute writes the Tink MAC tag of --in; circe mac verify checks a --tag against --in and exits 0, or 3 when it does not match.
status: stable
code:
  - internal/cli/mac.go
  - internal/cli/mac_test.go
tags: [cli, tink, primitives, mac]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T17:00:00Z }
---

# Intent
Scripts can authenticate data with a MAC keyset, and can check tags
that services produced.

# Behavior
- `circe mac compute --keyset <k>` runs `mac.New(h).ComputeMAC(data)` on
  `--in` and writes the tag, prefix included.
- `circe mac verify --keyset <k> --tag <path>` runs `VerifyMAC(tag, data)`
  with the data on `--in`. `--tag` is required and is read in
  `--encoding`. A tag that does not match is `circe: MAC verification failed`,
  exit 3. A match prints nothing and exits 0.
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `mac.New`, `tink.MAC`.
