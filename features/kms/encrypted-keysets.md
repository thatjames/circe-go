---
type: Feature
title: KMS-encrypted keysets
description: --master-key-uri encrypts and decrypts keysets with a key in GCP KMS, AWS KMS or HashiCorp Vault, via the tink-go KMS extensions.
status: draft
code:
  - internal/kms/**
tags: [kms, gcp, aws, vault]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T18:32:00Z }
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
- The KMS AEAD encrypts the keyset with Tink's envelope rules (the
  `EncryptedKeyset` proto), so any Tink library with the same KMS key
  can read it.
- KMS errors (permission, unknown key, network) are reported with the
  backend and key URI, never with credential values.
- Tests use fake AEADs in place of the KMS; CI needs no cloud access.

# Interfaces
- [Keyset I/O](/keyset/keyset-io.md) calls this package to get an AEAD
  for a URI.
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
