---
type: Feature
title: AEAD encrypt and decrypt
description: circe aead encrypt and decrypt run Tink's AEAD primitive on --in with an AEAD keyset and optional associated data.
status: draft
code:
  - internal/cli/aead.go
  - internal/cli/aead_test.go
tags: [cli, tink, primitives, aead]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T14:30:13Z }
---

# Intent
AEAD is Tink's default way to encrypt data. With this command a script
can encrypt and decrypt with the same keyset its services use.

# Behavior
- `circe aead encrypt --keyset <k>` runs `aead.New(h).Encrypt(plaintext, ad)`
  on the whole of `--in` and writes the ciphertext.
- `circe aead decrypt --keyset <k>` runs `Decrypt(ciphertext, ad)` and
  writes the plaintext. A wrong key, wrong associated data or a tampered
  ciphertext is `circe: decryption failed`, exit 3.
- `--associated-data <text>` or `--associated-data-file <path>` gives the
  associated data. The two flags are mutually exclusive (giving both is a
  usage error, exit 2). The default is empty.
- The input is read into memory whole, as Tink's AEAD API needs. Large
  files belong to [streaming AEAD](/primitives/streaming-aead.md).
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `aead.New`, `tink.AEAD`.
