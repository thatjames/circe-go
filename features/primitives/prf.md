---
type: Feature
title: PRF compute
description: circe prf compute --output-length <n> writes n bytes of the keyset's primary PRF over --in.
status: draft
code:
  - internal/cli/prf.go
  - internal/cli/prf_test.go
tags: [cli, tink, primitives, prf]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T11:20:00Z }
---

# Intent
Scripts and test fixtures can derive the same pseudorandom values, for
example identifiers, that a service derives with a PRF keyset.

# Behavior
- `circe prf compute --keyset <k> --output-length <n>` runs
  `prf.NewPRFSet(h)` and `ComputePrimaryPRF(input, n)` on `--in`, and
  writes `n` bytes. `--output-length` is required and must be at least 1.
- A length larger than the primary PRF allows (for example more than the
  hash size of HMAC-PRF) is an error that names the limit, exit 1.
- Only the primary PRF is used.
- Flags, encoding and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `prf.NewPRFSet`, `prf.Set`.

# Decisions
- **Primary PRF only** (2026-11-03, proposed): Tink's PRF set exposes
  every key by key ID, but scripts want the value a service computes, and
  a service uses the primary key. A `--key-id` flag can come later.
