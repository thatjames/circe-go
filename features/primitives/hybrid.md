---
type: Feature
title: Hybrid encrypt and decrypt
description: circe hybrid encrypt encrypts --in to a public hybrid keyset (or a private one) with optional context info; circe hybrid decrypt decrypts with the private keyset.
status: draft
code:
  - internal/cli/hybrid.go
  - internal/cli/hybrid_test.go
tags: [cli, tink, primitives, hybrid]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T11:20:00Z }
---

# Intent
Anyone who holds a public keyset can encrypt data for its owner (HPKE or
ECIES), for example to send a secret to a service.

# Behavior
- `circe hybrid encrypt --keyset <k>` runs
  `hybrid.NewHybridEncrypt(h).Encrypt(plaintext, contextInfo)` on `--in`
  and writes the ciphertext. It accepts a public keyset, or a private
  one, from which it uses `h.Public()`.
- `circe hybrid decrypt --keyset <private>` runs
  `hybrid.NewHybridDecrypt(h).Decrypt(ciphertext, contextInfo)`. A public
  keyset is an error saying that decryption needs the private keyset. A
  failure is `circe: decryption failed`, exit 1.
- `--context-info <text>` or `--context-info-file <path>` gives Tink's
  context info. The two flags are mutually exclusive (giving both is a
  usage error, exit 2). The default is empty.
- The input is read into memory whole.
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `hybrid.NewHybridEncrypt`, `hybrid.NewHybridDecrypt`.

# Decisions
- **`--context-info`, not `--associated-data`** (2026-11-03, proposed):
  the flag keeps Tink's name for the hybrid parameter.
