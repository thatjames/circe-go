---
type: Feature
title: Keyset lifecycle
description: create-keyset, add-key, promote-key, rotate-keyset, enable-key, disable-key, destroy-key and delete-key manage keys in a Tink keyset.
status: draft
code:
  - internal/cli/create_keyset*.go
  - internal/cli/add_key*.go
  - internal/cli/promote_key*.go
  - internal/cli/rotate_keyset*.go
  - internal/cli/enable_key*.go
  - internal/cli/disable_key*.go
  - internal/cli/destroy_key*.go
  - internal/cli/delete_key*.go
tags: [keyset, cli]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
The core of tinkey: create a keyset and manage its keys over time, using
tink-go's `keyset.Manager`.

# Behavior
- `create-keyset --key-template <NAME>` writes a new keyset with one
  enabled primary key from the [template](/keyset/key-templates.md).
- `add-key --key-template <NAME>` reads a keyset, adds a new enabled,
  non-primary key and writes the keyset.
- `promote-key --key-id <ID>` makes an enabled key the primary.
- `rotate-keyset --key-template <NAME>` adds a key and makes it primary
  (add-key plus promote-key in one step).
- `enable-key`, `disable-key`, `destroy-key`, `delete-key --key-id <ID>`
  change one key's status or remove it, with tink-go's rules: the primary
  cannot be disabled, destroyed or deleted; destroying wipes the key
  material and keeps the entry.
- An unknown `--key-id` is an error; the output is not written.
- Each command reads and writes through [keyset I/O](/keyset/keyset-io.md)
  with the shared flags of the [command surface](/cli/command-surface.md),
  so they work on cleartext and encrypted keysets alike.
- Each command has tests that run it on synthetic keysets and check the
  resulting keyset.

# Interfaces
- `keyset.Manager` of tink-go v2.

# Decisions
- **All tinkey commands in scope** (user, 2026-10-22), built
  incrementally.
- `rotate-keyset` is kept although Java tinkey deprecated it: it is a
  convenience, and the CLI surface is free. *Open:* the user did not
  confirm keeping it at the import (2026-10-22); decide before it is
  built.
