---
type: Feature
title: Key templates
description: Maps Tink key template names (AES256_GCM, ECDSA_P256, ...) to tink-go templates for --key-template, and lists them.
status: draft
code:
  - internal/templates/**
tags: [keyset, tink, templates]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Users name the kind of key they want with Tink's template names; tink-go
has template functions but no lookup by name, so circe keeps the table.

# Behavior
- `--key-template <NAME>` resolves `NAME` to a tink-go
  `*tinkpb.KeyTemplate`. Names are Tink's template names, upper case
  (for example `AES128_GCM`, `AES256_GCM`, `AES256_GCM_SIV`,
  `CHACHA20_POLY1305`, `XCHACHA20_POLY1305`, `AES256_SIV`,
  `HMAC_SHA256_256BITTAG`, `ECDSA_P256`, `ED25519`,
  `DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM`).
- The table covers every template tink-go exposes for the primitives it
  supports (AEAD, deterministic AEAD, streaming AEAD, MAC, PRF,
  signature, hybrid, JWT).
- An unknown name is an error that suggests `circe list-key-templates`.
- `circe list-key-templates` prints the names, one per line, sorted.
- A test asserts every name produces a template tink-go can create a
  keyset from.

# Interfaces
- Used by `create-keyset` and `add-key` ([lifecycle](/keyset/lifecycle.md)).

# Decisions
- **Tink key templates** (user, 2026-10-22): names and parameters come
  from Tink, not from Java tinkey's list. Where the two differ, Tink's
  name wins.
