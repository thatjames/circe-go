---
type: Feature
title: Keyset I/O
description: Reads and writes cleartext Tink keysets in Tink's JSON and binary formats through tink-go only; the hook for KMS-encrypted keysets.
status: stable
code:
  - internal/keysetio/**
tags: [keyset, tink]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:06:05Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
A keyset circe writes must be readable by any Tink library, and any Tink
keyset must be readable by circe. Tink's formats are the contract.

# Behavior
- Formats: `json` (Tink's JSON keyset, `keyset.NewJSONReader/Writer`)
  and `binary` (the serialized `Keyset` proto).
- Keysets are read and written in cleartext with tink-go's
  `insecurecleartextkeyset`.
- When a command is given `--master-key-uri`, reading and writing are
  encrypted instead; that path is described and built by
  [KMS-encrypted keysets](/kms/encrypted-keysets.md).
- A keyset that does not parse in the given format is an error naming
  the input (its path, or `stdin`) and the format; nothing is written.
- An empty input is an error too.
- JSON keysets are written as tink-go writes them (one line), followed by
  a newline.
- Round trip: reading a keyset and writing it in the same format yields
  the same keyset (same key IDs, primary, statuses and key material).

# Interfaces
- Used by every keyset command: [create](/keyset/create-keyset.md),
  [key management](/keyset/key-management.md),
  [public keyset](/keyset/public-keyset.md),
  [inspect and convert](/keyset/inspect-convert.md).
- Flags: [command surface](/cli/command-surface.md).
- `keysetio.Format` (`json`, `binary`; a `pflag.Value`),
  `keysetio.Read(r io.Reader, name string, f Format) (*keyset.Handle, error)`
  and `keysetio.Write(w io.Writer, h *keyset.Handle, f Format) error`.

# Decisions
- **Tink keyset formats** (user, 2026-10-22): read and write Tink's JSON
  and binary keysets exactly; never a circe-specific format.
- **Build on tink-go** (user, 2026-10-22): no hand-rolled crypto or
  proto handling where tink-go has an API.
- Encrypted I/O is part of the KMS feature, so this one is fully built
  by the CLI foundation (user, 2026-10-22).
