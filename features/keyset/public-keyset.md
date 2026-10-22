---
type: Feature
title: Create public keyset
description: create-public-keyset reads a private signature, hybrid or JWT signature keyset and writes its public keyset in cleartext.
status: draft
code:
  - internal/cli/create_public_keyset*.go
tags: [keyset, cli]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:10:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Hand out the public half of asymmetric keys without the private keys.

# Behavior
- `create-public-keyset` reads a private keyset of an asymmetric
  primitive (signature, hybrid, JWT signature) and writes its public
  keyset in cleartext (`keyset.Handle.Public`).
- A keyset without private keys is an error; nothing is written.
- Input and output go through [keyset I/O](/keyset/keyset-io.md) with
  the flags of the [command surface](/cli/command-surface.md).
- Tests on synthetic signature, hybrid and JWT keysets check that the
  output holds only public keys and verifies what the private keyset
  signs.

# Decisions
- **All tinkey commands in scope** (user, 2026-10-22).
