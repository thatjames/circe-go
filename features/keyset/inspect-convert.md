---
type: Feature
title: Inspect and convert keysets
description: list-keyset shows a keyset's key info without key material; convert-keyset rewrites a keyset between Tink's JSON and binary formats.
status: draft
code:
  - internal/cli/list_keyset*.go
  - internal/cli/list_key_templates*.go
  - internal/cli/convert_keyset*.go
tags: [keyset, cli]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:10:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Look at a keyset safely and move it between Tink's formats.

# Behavior
- `list-keyset` prints the keyset's `KeysetInfo` (primary key ID; per
  key: ID, type URL, status, output prefix type). It never prints key
  material.
- `list-key-templates`: see [key templates](/keyset/key-templates.md).
- `convert-keyset` reads with the `--in*` flags and writes with the
  `--out*` flags, so `--in-format` and `--out-format` may differ.
- Moving a keyset between cleartext and KMS encryption, or between KMS
  keys, is part of [KMS-encrypted keysets](/kms/encrypted-keysets.md).
- Each command has tests on synthetic keysets.

# Interfaces
- [Keyset I/O](/keyset/keyset-io.md), [command surface](/cli/command-surface.md).

# Decisions
- **All tinkey commands in scope** (user, 2026-10-22).
- Encrypted conversion lives in the KMS feature, so this one is fully
  built by the CLI foundation (user, 2026-10-22).
