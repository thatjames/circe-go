---
type: Feature
title: Primitive commands
description: Shared shape of circe's per-primitive commands (aead, daead, streaming-aead, mac, prf, signature, hybrid, JWT) - --keyset and its KMS flags, data on --in/--out, --encoding raw|base64|hex, and wrong-keyset and verification errors.
status: draft
code:
  - internal/cli/primitive_flags*.go
  - internal/cli/encoding*.go
tags: [cli, tink, primitives]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T11:20:00Z }
---

# Intent
circe uses a Tink keyset for its operations as well as managing it, so
that shell scripts, interop debugging and test fixtures can encrypt,
decrypt, sign, verify and compute MACs and PRFs from the command line.
There is one command group per Tink primitive, and every group works the
same way.

# Behavior
- One command group per Tink primitive, each with its operations:
  [aead](/primitives/aead.md), [daead](/primitives/deterministic-aead.md),
  [streaming-aead](/primitives/streaming-aead.md), [mac](/primitives/mac.md),
  [prf](/primitives/prf.md), [signature](/primitives/signature.md),
  [hybrid](/primitives/hybrid.md) and [JWT](/primitives/jwt.md). A group
  given no operation prints its help and exits 0.
- `--keyset <path>` names the keyset and is required, because the data
  comes on stdin. `--keyset-format json|binary` defaults to `json`.
  `--master-key-uri` and `--credential` read a KMS-encrypted keyset
  ([encrypted keysets](/kms/encrypted-keysets.md)). The keyset is only
  read, never written.
- The data comes from `--in <path>` (default: stdin) and the result goes
  to `--out <path>` (default: stdout). `--out` follows the
  [command surface](/cli/command-surface.md) rules: it never overwrites a
  file, a file it creates has mode 0600 (decrypted data is secret), and
  nothing is written when the command fails.
- `--encoding raw|base64|hex` (default `raw`) is the encoding of the bytes
  Tink produces: ciphertexts, MAC tags, signatures and PRF output. It
  applies to them both as output (encrypt, compute, sign) and as input
  (decrypt's `--in`, verify's `--tag` or `--signature`). Plaintext and
  signed or MACed data are always raw bytes.
  - `base64` is RFC 4648 standard base64 with padding. `hex` is lower
    case. Encoded output ends in a newline. Encoded input may have
    surrounding whitespace; input that does not decode is an error that
    names the input and the encoding.
- The primary key encrypts, signs and computes. Decrypt and verify try
  every enabled key, matched by output prefix, as Tink does.
- A keyset that holds no key of the group's primitive is an error naming
  the keyset, the primitive the command expected and the key type found,
  for example `circe: k.json is not an AEAD keyset (type.googleapis.com/google.crypto.tink.HmacKey)`.
  The exit code is 1.
- A failed decryption or verification prints `circe: <operation> failed`
  on stderr and exits 1. Nothing goes to stdout and no `--out` file is
  written. A successful verify prints nothing and exits 0, so scripts can
  write `if circe signature verify ...`.
- No error message holds key material, plaintext or associated data.
- Tests run every operation with synthetic keysets made from
  [key templates](/keyset/key-templates.md): round trips, every encoding,
  wrong-primitive keysets, tampered input. They also check circe's output
  with tink-go's primitives directly, which proves the output is Tink's
  wire format.

# Interfaces
- Keysets are read through [keyset I/O](/keyset/keyset-io.md)
  (`keysetio.Read`, `keysetio.ReadEncrypted`).
- `--out` is written through the same write-on-success path as the
  keyset commands ([command surface](/cli/command-surface.md)).
- The primitive groups are listed under the command surface's
  subcommands.

# Decisions
- **Grouped by primitive** (user, 2026-11-03): commands are
  `circe <primitive> <operation>`, for example `circe aead encrypt`. The
  primitive is never guessed from the keyset.
- **Raw bytes by default, `--encoding` for text** (user, 2026-11-03).
- **All tink-go primitives** (user, 2026-11-03): AEAD, deterministic AEAD,
  streaming AEAD, MAC, PRF, signature, hybrid and JWT.
- **`--keyset`, not `--in`, for the keyset** (2026-11-03, proposed): on
  these commands `--in` is the data, so the keyset needs a flag of its
  own.
