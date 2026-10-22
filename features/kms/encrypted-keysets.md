---
type: Feature
title: KMS-encrypted keysets
description: --master-key-uri reads and writes KMS-encrypted keysets on every command (GCP KMS, AWS KMS, Vault), and convert-keyset moves keysets between cleartext and KMS keys.
status: stable
code:
  - internal/kms/**
  - internal/keysetio/encrypted*.go
  - internal/cli/kms_flags*.go
tags: [kms, gcp, aws, vault]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-22T19:45:23Z }
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
  - Vault: token auth only. `--credential` is a file holding the token
    (surrounding whitespace ignored), else `VAULT_TOKEN`; no token is an
    error. The server is the URI's host over HTTPS
    (`hcvault://vault.example.com:8200/transit/keys/<name>`); a URI with
    no host (`hcvault:///transit/keys/<name>`) uses `VAULT_ADDR`. TLS
    settings come from Vault's own environment (`VAULT_CACERT`, ...).
- A backend client is made only when a command is given a URI for it;
  nothing contacts a KMS before the first encrypt or decrypt.
- Every keyset command accepts `--master-key-uri` and `--credential`.
  With them, [keyset I/O](/keyset/keyset-io.md) reads the input as an
  `EncryptedKeyset` (JSON or binary, per `--in-format`) decrypted with the
  KMS AEAD (`keyset.Read`), and writes the output encrypted with the same
  AEAD (`Handle.Write`), so any Tink library with the same KMS key can
  read it. The associated data is empty, as in Tink's `keyset.Read`.
  - `create-keyset` only writes: the new keyset is written encrypted.
  - `list-keyset` only reads.
  - `create-public-keyset` reads the private keyset encrypted and writes
    the public keyset in cleartext, as [public keyset](/keyset/public-keyset.md)
    says: it holds no secrets.
- Reading an encrypted keyset without `--master-key-uri`, or a cleartext
  one with it, fails to parse and is an error; nothing is written.
- `convert-keyset` takes `--new-master-key-uri` and `--new-credential`
  for the output side: input and output may use different KMS keys or
  backends. Omitting `--new-master-key-uri` writes cleartext; giving it
  with a cleartext input encrypts the keyset.
- A keyset that fails to decrypt is an error naming the input and the
  key URI; nothing is written.
- KMS errors (permission, unknown key, network) are reported with the
  backend and key URI, never with credential values.
- Tests use fake AEADs in place of the KMS; CI needs no cloud access.
  The backend dispatch is tested offline: clients are built, never used.

# Interfaces
- `kms.AEAD(ctx, uri, credential string) (tink.AEAD, error)` returns the
  AEAD for a key URI; `keysetio.ReadEncrypted(r, name, format, aead, uri)`
  and `keysetio.WriteEncrypted(w, h, format, aead)` live beside the
  cleartext [keyset I/O](/keyset/keyset-io.md) in
  `internal/keysetio/encrypted*.go`.
- `internal/cli/kms_flags.go` adds the flags and hands the AEAD to the
  shared `--in`/`--out` flags; one AEAD serves both sides of a command.
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
