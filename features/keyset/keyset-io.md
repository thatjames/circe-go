---
type: Feature
title: Keyset I/O
description: Reads and writes Tink keysets in Tink's JSON and binary formats, cleartext or KMS-encrypted, through tink-go only.
status: draft
code:
  - internal/keysetio/**
tags: [keyset, tink]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
A keyset circe writes must be readable by any Tink library, and any Tink
keyset must be readable by circe. Tink's formats are the contract.

# Behavior
- Formats: `json` (Tink's JSON keyset, `keyset.NewJSONReader/Writer`)
  and `binary` (the serialized `Keyset` / `EncryptedKeyset` proto).
- Without `--master-key-uri`, a keyset is read and written in cleartext
  with tink-go's `insecurecleartextkeyset`.
- With `--master-key-uri`, the input is an `EncryptedKeyset` decrypted
  with the KMS AEAD, and the output is encrypted with it
  ([encrypted keysets](/kms/encrypted-keysets.md)).
- A keyset that does not parse in the given format, or that fails to
  decrypt, is an error naming the input and the format; nothing is
  written.
- Round trip: reading a keyset and writing it in the same format yields
  the same keyset (same key IDs, primary, statuses and key material).

# Interfaces
- Used by every command in [lifecycle](/keyset/lifecycle.md) and
  [inspect and convert](/keyset/inspect-convert.md).
- Flags: [command surface](/cli/command-surface.md).

# Decisions
- **Tink keyset formats** (user, 2026-10-22): read and write Tink's JSON
  and binary keysets exactly; never a circe-specific format.
- **Build on tink-go** (user, 2026-10-22): no hand-rolled crypto or
  proto handling where tink-go has an API.
