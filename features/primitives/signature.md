---
type: Feature
title: Signature sign and verify
description: circe signature sign signs --in with a private signature keyset; circe signature verify checks a --signature against --in with the public (or private) keyset and exits 0 or 1.
status: draft
code:
  - internal/cli/signature.go
  - internal/cli/signature_test.go
tags: [cli, tink, primitives, signature]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T11:20:00Z }
---

# Intent
Scripts can sign release artifacts and configs, and can verify
signatures made by any Tink library.

# Behavior
- `circe signature sign --keyset <private>` runs
  `signature.NewSigner(h).Sign(data)` on `--in` and writes the signature.
  A public keyset is an error saying that signing needs the private
  keyset.
- `circe signature verify --keyset <k> --signature <path>` runs
  `signature.NewVerifier(h).Verify(sig, data)` with the data on `--in`.
  `--signature` is required and is read in `--encoding`. A signature that
  does not verify is `circe: signature verification failed`, exit 1. A
  valid one prints nothing and exits 0.
- verify accepts a private keyset as well as a public one: from a private
  keyset it uses `h.Public()`, as
  [create-public-keyset](/keyset/public-keyset.md) does.
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `signature.NewSigner`, `signature.NewVerifier`.
