---
type: Feature
title: Key management
description: add-key, promote-key, enable-key, disable-key and delete-key change the keys of an existing Tink keyset; destroy-key waits on a decision; no rotate-keyset.
status: stable
code:
  - internal/cli/add_key*.go
  - internal/cli/promote_key*.go
  - internal/cli/enable_key*.go
  - internal/cli/disable_key*.go
  - internal/cli/delete_key*.go
  - internal/cli/key_commands*.go
tags: [keyset, cli]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:10:39Z }
verified: { by: human:slimjim, at: 2026-10-22T19:00:00Z }
---

# Intent
Manage a keyset's keys over time with tink-go's `keyset.Manager`, one
explicit step per command.

# Behavior
- `add-key --key-template <NAME>` reads a keyset, adds a new enabled,
  non-primary key from the [template](/keyset/key-templates.md) and
  writes the keyset.
- `promote-key --key-id <ID>` makes an enabled key the primary; a
  disabled key is an error.
- `enable-key`, `disable-key`, `delete-key --key-id <ID>` change one
  key's status or remove it, with tink-go's rules: the primary cannot be
  disabled or deleted. Enabling an enabled key, or disabling a disabled
  one, is not an error.
- `destroy-key` is not offered yet: see Decisions.
- A missing `--key-id` is a usage error. An unknown one is an error
  (`no key with ID <ID> in the keyset`); the output is not written.
- Errors from tink-go are shown without its `keyset.Manager: ` prefix.
- The shared `--key-id` handling lives in `key_commands.go`.
- Each command reads and writes through [keyset I/O](/keyset/keyset-io.md)
  with the shared flags of the [command surface](/cli/command-surface.md).
- Each command has tests that run it on synthetic keysets and check the
  resulting keyset.

# Interfaces
- `keyset.Manager` of tink-go v2.
- Creating the keyset: [create keyset](/keyset/create-keyset.md).

# Decisions
- **All tinkey commands in scope** (user, 2026-10-22), built
  incrementally.
- **No `rotate-keyset`** (user, 2026-10-22): rotating is `add-key` then
  `promote-key`, explicitly. A combined command blurs which key is
  primary and can confuse users; anyone who wants it can alias the two.
- **`destroy-key` waits on the user** (2026-10-22): Tink destroys a key
  by setting `DESTROYED` and removing its key data, but tink-go v2.9.0
  rejects such a keyset on read (`key <ID> has no key data`) and its
  `keyset.Manager` has no destroy. A circe-destroyed keyset would be
  readable by Java, C++ and Python Tink, but not by tink-go or circe
  itself. Until the user decides, `destroy-key` is left out; `delete-key`
  removes a key entirely.
