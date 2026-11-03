---
type: Feature
title: Signature sign and verify
description: circe signature sign signs --in with a private signature keyset; circe signature verify checks a --signature against --in with the public keyset and exits 0, or 3 when it does not verify.
status: stable
code:
  - internal/cli/signature.go
  - internal/cli/signature_test.go
tags: [cli, tink, primitives, signature]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T17:00:00Z }
---

# Intent
Scripts can sign release artifacts and configs, and can verify
signatures made by any Tink library.

# Behavior
- `circe signature sign --keyset <private>` runs
  `signature.NewSigner(h).Sign(data)` on `--in` and writes the signature.
  A public keyset is an error saying that signing needs the private
  keyset.
- `circe signature verify --keyset <public> --signature <path>` runs
  `signature.NewVerifier(h).Verify(sig, data)` with the data on `--in`.
  `--signature` is required and is read in `--encoding`. A signature that
  does not verify is `circe: signature verification failed`, exit 3. A
  valid one prints nothing and exits 0.
- verify refuses a private keyset with an error that says to make the
  public keyset with [create-public-keyset](/keyset/public-keyset.md)
  first, exit 1.
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `signature.NewSigner`, `signature.NewVerifier`.
