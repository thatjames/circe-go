---
type: Feature
title: Inspect and convert keysets
description: list-keyset shows a keyset's key info without key material, convert-keyset changes format or encryption, create-public-keyset extracts the public keyset.
status: draft
code:
  - internal/cli/list_keyset*.go
  - internal/cli/list_key_templates*.go
  - internal/cli/convert_keyset*.go
  - internal/cli/create_public_keyset*.go
tags: [keyset, cli]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Look at a keyset safely, move it between formats and between cleartext
and KMS encryption, and hand out the public half of asymmetric keys.

# Behavior
- `list-keyset` prints the keyset's `KeysetInfo` (primary key ID; per
  key: ID, type URL, status, output prefix type). It never prints key
  material.
- `list-key-templates`: see [key templates](/keyset/key-templates.md).
- `convert-keyset` reads with the `--in*` flags and writes with the
  `--out*` flags. Input and output may differ in format and in
  encryption: `--master-key-uri` for the input and `--new-master-key-uri`
  for the output, with `--credential` and `--new-credential` for their
  credentials; omitting `--new-master-key-uri` writes cleartext.
- `create-public-keyset` reads a private keyset of an asymmetric
  primitive (signature, hybrid, JWT signature) and writes its public
  keyset in cleartext. A keyset without private keys is an error.
- Each command has tests on synthetic keysets.

# Interfaces
- [Keyset I/O](/keyset/keyset-io.md), [command surface](/cli/command-surface.md),
  [encrypted keysets](/kms/encrypted-keysets.md).

# Decisions
- **All tinkey commands in scope** (user, 2026-10-22).
- **`--new-master-key-uri` and `--new-credential`** for the output side
  of `convert-keyset` (user, 2026-10-22).
