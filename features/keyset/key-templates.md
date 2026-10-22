---
type: Feature
title: Key templates
description: Maps Tink key template names (AES256_GCM, ECDSA_P256, ...) to tink-go templates for --key-template, and lists them.
status: stable
code:
  - internal/templates/**
tags: [keyset, tink, templates]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:06:05Z }
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
- Each name is the one Tink registers for those parameters (the
  `namedParameters` of tink-java's key managers, which `tinkey
  list-key-templates` prints). Where Tink registers two names for the
  same parameters (`ECDSA_P384` and `ECDSA_P384_SHA512`, `AES_CMAC` and
  `AES256_CMAC`, `RSA_SSA_PSS_3072_SHA256_F4` and
  `RSA_SSA_PSS_3072_SHA256_SHA256_32_F4`, ...), both resolve.
- Tink's `_RAW` name of a template tink-go only has with the `TINK`
  prefix (for example `CHACHA20_POLY1305_RAW`, `HMAC_SHA256_256BITTAG_RAW`,
  `AES256_SIV_RAW`) resolves to that template with output prefix `RAW`.
- tink-go templates Tink registers no name for take Tink's naming
  pattern: the DER, no-prefix ECDSA templates are `ECDSA_P256_NO_PREFIX`,
  `ECDSA_P384_NO_PREFIX`, `ECDSA_P384_SHA384_NO_PREFIX`,
  `ECDSA_P521_NO_PREFIX`, and the prefixed 160-bit-nonce XAES-GCM is
  `XAES_256_GCM_160_BIT_NONCE`.
- Tink names with no tink-go template (AES-EAX, the IEEE P1363 ECDSA
  templates with a Tink prefix, compressed-point ECIES, the P384 and P521
  HPKE templates, ML-DSA, SLH-DSA) are not offered.
- Streaming AEAD and PRF keys never carry a prefix in Tink, so their
  templates (no `_RAW` in the name) give `RAW` keys.
- Names match exactly (upper case). An unknown name is an error that
  suggests `circe list-key-templates`.
- `circe list-key-templates` prints the names, one per line, sorted.
- A test asserts every name produces a template tink-go can create a
  keyset from, with the output prefix its name says.

# Interfaces
- Used by [create-keyset](/keyset/create-keyset.md) and
  [add-key](/keyset/key-management.md).

# Decisions
- **Tink key templates** (user, 2026-10-22): names and parameters come
  from Tink, not from Java tinkey's list. Where the two differ, Tink's
  name wins.
- **Raw variants by output prefix** (2026-10-22): tink-go has no API for
  the `_RAW` templates it lacks, so circe copies the tink-go template and
  sets only `OutputPrefixType` to `RAW`; every other parameter stays
  tink-go's.
- `ED25519WithRawOutput`, Tink's mixed-case alias of `ED25519_RAW`, is
  left out so every name is upper case.
