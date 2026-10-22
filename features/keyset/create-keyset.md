---
type: Feature
title: Create keyset
description: create-keyset writes a new Tink keyset with one enabled primary key from a named Tink key template.
status: stable
code:
  - internal/cli/create_keyset*.go
tags: [keyset, cli]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:06:05Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
The first step of every keyset's life: make one from a Tink template,
using tink-go's `keyset.Manager`.

# Behavior
- `create-keyset --key-template <NAME>` writes a new keyset with one
  enabled primary key from the [template](/keyset/key-templates.md).
- A missing `--key-template` is a usage error.
- The keyset is written through [keyset I/O](/keyset/keyset-io.md) with
  the `--out*` flags of the [command surface](/cli/command-surface.md).
- Tests create a keyset from several templates and check that tink-go
  loads it with one enabled primary key.

# Interfaces
- `keyset.NewHandle` / `keyset.Manager` of tink-go v2.
- Changing the keyset afterwards: [key management](/keyset/key-management.md).

# Decisions
- **All tinkey commands in scope** (user, 2026-10-22), built
  incrementally; split from key management so each feature is built by
  one plan (user, 2026-10-22).
