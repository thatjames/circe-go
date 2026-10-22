---
type: Feature
title: KMS-encrypted keysets
description: --master-key-uri reads and writes KMS-encrypted keysets on every command (GCP KMS, AWS KMS, Vault), and convert-keyset moves keysets between cleartext and KMS keys.
status: draft
code:
  - internal/kms/**
  - internal/keysetio/encrypted*.go
tags: [kms, gcp, aws, vault]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:10:00Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Keysets at rest should be encrypted by a key the user keeps in a KMS, as
Tink recommends, on the three backends the user works with.

# Behavior
- `--master-key-uri` selects the backend by prefix:
  - `gcp-kms://...` via `tink-go-gcpkms/v2`;
  - `aws-kms://...` via `tink-go-awskms/v2`;
  - `hcvault://...` via `tink-go-hcvault/v2`.
  Any other prefix is an error listing the supported ones.
- `--credential <path>` gives the backend's credentials file. Without
  it, each backend uses its default credential chain:
  - GCP: Application Default Credentials (`GOOGLE_APPLICATION_CREDENTIALS`);
  - AWS: the SDK default chain (`AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`,
    `AWS_PROFILE`, shared config);
  - Vault: token auth only, `VAULT_TOKEN` and `VAULT_ADDR`;
    `--credential` is a file holding the token.
- Every keyset command accepts `--master-key-uri` and `--credential`.
  With them, [keyset I/O](/keyset/keyset-io.md) reads the input as an
  `EncryptedKeyset` (JSON or binary, per `--in-format`) decrypted with the
  KMS AEAD (`keyset.Read`), and writes the output encrypted with the same
  AEAD (`Handle.Write`), so any Tink library with the same KMS key can
  read it.
- `convert-keyset` takes `--new-master-key-uri` and `--new-credential`
  for the output side: input and output may use different KMS keys or
  backends. Omitting `--new-master-key-uri` writes cleartext; giving it
  with a cleartext input encrypts the keyset.
- A keyset that fails to decrypt is an error naming the input and the
  key URI; nothing is written.
- KMS errors (permission, unknown key, network) are reported with the
  backend and key URI, never with credential values.
- Tests use fake AEADs in place of the KMS; CI needs no cloud access.

# Interfaces
- [Keyset I/O](/keyset/keyset-io.md) calls this package to get an AEAD
  for a URI; the encrypted read/write lives beside it in
  `internal/keysetio/encrypted*.go`.
- [Inspect and convert](/keyset/inspect-convert.md): `convert-keyset`.
- Flags on the [command surface](/cli/command-surface.md).

# Decisions
- **Backends: GCP KMS, AWS KMS, HashiCorp Vault** (user, 2026-10-22).
- Credentials are never logged or printed; only key names of
  environment variables appear in docs and errors.
- **Vault: token only in v1** (user, 2026-10-22). AppRole or Kubernetes
  auth is later work.
- **One credential per side** (user, 2026-10-22): `--credential` for the
  input, `--new-credential` for the output of `convert-keyset`, so it can
  move a keyset between backends.
- Encrypted I/O and encrypted conversion belong to this feature, so the
  KMS plan builds it whole (user, 2026-10-22).
