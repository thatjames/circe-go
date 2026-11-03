---
type: Feature
title: JWT MAC and signature
description: circe signs or MACs a JWT claims JSON into a compact JWT with a Tink JWT keyset, and verifies a compact JWT against Tink validator flags, printing its verified claims.
status: draft
code:
  - internal/cli/jwt*.go
tags: [cli, tink, primitives, jwt]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T15:20:00Z }
---

# Intent
Test fixtures and debugging need real tokens: a JWT a service will
accept, and a check of why a service rejects one.

# Behavior
- Commands, one group per Tink JWT primitive:
  - JWT MAC keysets: `circe jwt-mac compute` runs
    `jwt.NewMAC(h).ComputeMACAndEncode(raw)`, and `circe jwt-mac verify`
    runs `VerifyMACAndDecode(compact, validator)`.
  - JWT signature keysets: `circe jwt-signature sign` runs
    `jwt.NewSigner(h).SignAndEncode(raw)` with the private keyset, and
    `circe jwt-signature verify` runs
    `jwt.NewVerifier(h).VerifyAndDecode(compact, validator)` with the
    public keyset. As in [signature](/primitives/signature.md), verify
    refuses a private keyset and sign refuses a public one, exit 1.
- Sign and compute read a claims JSON object on `--in` and write the
  compact JWT followed by a newline. `--type-header <typ>` sets the `typ`
  header.
  - The registered claims must have the types JWT gives them (`exp`,
    `nbf` and `iat` are numbers, `aud` is a string or an array of
    strings). Any other value is an error naming the claim.
  - Claims without `exp` are refused unless `--without-expiration` is
    given, and claims with `exp` are refused when it is given: the same
    rules as Tink's `RawJWTOptions`. Both are errors, exit 1.
  - Input that is not a JSON object is an error, exit 1.
- Verify reads a compact JWT on `--in` (surrounding whitespace ignored).
  It builds a `jwt.Validator` from `--expected-type-header`,
  `--expected-issuer`, `--expected-audience`, `--ignore-type-header`,
  `--ignore-issuer`, `--ignore-audiences`, `--allow-missing-expiration`,
  `--expect-issued-in-the-past`, `--clock-skew <duration>` (at most 10m,
  Tink's limit) and `--fixed-now <RFC 3339 time>` (for fixtures and
  debugging).
  - Flags Tink refuses are usage errors, exit 2: a clock skew over 10m, a
    `--fixed-now` that is not RFC 3339, or an expected value given with
    its ignore flag (for example `--expected-issuer` with `--ignore-issuer`).
  - On success it writes the verified claims as JSON
    (`VerifiedJWT.JSONPayload`) and a newline, and exits 0.
  - On failure it prints `circe: JWT verification failed: <Tink's reason>`
    and exits 3. The reason says, for example, that the token expired or
    that the issuer did not match, and never contains the token. tink-go
    quotes a segment of the token when that segment is not base64url, so
    a reason that holds any segment of the token is replaced by
    `the token is not a valid JWS compact serialization`.
- JWTs are text already, so these commands take no `--encoding`.
- Other flags and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `jwt.NewMAC`, `jwt.NewSigner`, `jwt.NewVerifier`,
  `jwt.NewRawJWTFromJSON`, `jwt.NewValidator`.

# Decisions
- **`jwt-mac` and `jwt-signature`** (user, 2026-11-03): one command group
  per Tink JWT primitive, the rule every other group follows. One `jwt`
  group that guessed the primitive from the keyset was rejected.
- **Claims through `jwt.NewRawJWTFromJSON`** (user, 2026-11-03): it
  keeps the claims exactly as given and checks the registered claim
  types, as Tink's own decoder does. Its doc comment points users to
  `NewRawJWT`, which would mean rebuilding `RawJWTOptions` from the JSON.
  It does not apply the `exp` rules, so circe checks them itself.
- JWK set import and export (`jwt.JWKSetFromPublicKeysetHandle`) is out
  of scope. It would be a keyset command and gets its own plan.
