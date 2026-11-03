---
type: Feature
title: Deterministic AEAD encrypt and decrypt
description: circe daead encrypt and decrypt run Tink's deterministic AEAD on --in with a DAEAD keyset; equal input and associated data give an equal ciphertext.
status: draft
code:
  - internal/cli/daead.go
  - internal/cli/daead_test.go
tags: [cli, tink, primitives, daead]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T11:20:00Z }
---

# Intent
Deterministic encryption (AES-SIV) lets a script compute the ciphertext
a service stores or looks up, for example to match an encrypted column.

# Behavior
- `circe daead encrypt --keyset <k>` runs
  `daead.New(h).EncryptDeterministically(plaintext, ad)` and writes the
  ciphertext. The same plaintext, associated data and primary key always
  give the same ciphertext.
- `circe daead decrypt --keyset <k>` runs `DecryptDeterministically`. A
  failure is `circe: decryption failed`, exit 1.
- `--associated-data` and `--associated-data-file` work as in
  [AEAD](/primitives/aead.md). The input is read into memory whole.
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `daead.New`, `tink.DeterministicAEAD`.
