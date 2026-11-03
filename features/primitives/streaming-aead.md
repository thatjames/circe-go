---
type: Feature
title: Streaming AEAD encrypt and decrypt
description: circe streaming-aead encrypt and decrypt stream --in to --out in constant memory with Tink's streaming AEAD; raw bytes only, and a failed decrypt leaves no --out file.
status: draft
code:
  - internal/cli/streaming_aead.go
  - internal/cli/streaming_aead_test.go
tags: [cli, tink, primitives, streaming-aead]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-03T11:20:00Z }
---

# Intent
Streaming AEAD encrypts files of any size, such as backups and
artifacts, without holding them in memory.

# Behavior
- `circe streaming-aead encrypt --keyset <k>` copies `--in` through
  `streamingaead.New(h).NewEncryptingWriter(w, ad)` to the output.
- `circe streaming-aead decrypt --keyset <k>` copies `--in` through
  `NewDecryptingReader(r, ad)` to the output.
- Memory use does not grow with the input size.
- `--associated-data` and `--associated-data-file` work as in
  [AEAD](/primitives/aead.md).
- No `--encoding`. Streams are raw bytes, and a user who wants text pipes
  the output through `base64`.
- With `--out`, the output goes to a temporary file in the same
  directory, which is renamed to `--out` only when the command succeeds
  and is removed when it fails. This keeps the write-on-success rule
  without buffering in memory. On stdout, plaintext that was already
  written before an authentication failure cannot be taken back: the
  command still exits 1 with `circe: decryption failed`, and the README
  tells users to decrypt to a file when that matters.
- Flags and errors: [primitive commands](/primitives/primitive-commands.md).

# Interfaces
- tink-go `streamingaead.New`, `tink.StreamingAEAD`.
