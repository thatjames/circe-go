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
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-26T09:07:48Z }
verified: { by: human:slimjim, at: 2026-10-22T18:40:00Z }
---

# Intent
Keysets at rest should be encrypted by a key the user keeps in a KMS, as
Tink recommends, on the three backends the user works with.

# Behavior
- `--master-key-uri` selects the backend by prefix:
  - `gcp-kms://...` via `tink-go-gcpkms/v2`;
  - `aws-kms://...` via `tink-go-awskms/v3` (aws-sdk-go-v2);
  - `hcvault://...` via `tink-go-hcvault/v2`.
  Any other prefix is an error listing the supported ones.
- `--credential <path>` gives the backend's credentials file. Without
  it, each backend uses its default credential chain:
  - GCP: Application Default Credentials (`GOOGLE_APPLICATION_CREDENTIALS`);
  - AWS: the aws-sdk-go-v2 default chain: `AWS_ACCESS_KEY_ID`/
    `AWS_SECRET_ACCESS_KEY`, then the profile named by `AWS_PROFILE`
    (else `default`) from `~/.aws/config` and `~/.aws/credentials`,
    including SSO, assume-role and `credential_process` profiles, then
    container and instance roles. The region always comes from the key
    URI. With `--credential`, a CSV file (as the IAM console downloads
    it) gives static keys; any other file is read as a shared
    credentials/config INI file (`[name]` or `[profile name]` sections)
    with the profile named by `AWS_PROFILE`, else `default`; that profile
    wins over `AWS_ACCESS_KEY_ID`, and `~/.aws` is not read. A profile missing from the file is an error naming
    the profile and the file. `--new-credential` follows the same rules;
    one `AWS_PROFILE` serves both sides;
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
- No backend SDK logs. `GOOGLE_SDK_GO_LOGGING_LEVEL` (any level) prints
  nothing: the GCP client gets a logger that discards everything, so
  tokens, the service-account assertion and the keyset plaintext of
  KMS Encrypt/Decrypt never reach stderr. AWS logs only with a
  `ClientLogMode` set in code, Vault only with a `Logger` set in code;
  circe sets neither.
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
- **GCP SDK logging is discarded, not redacted** (user, 2026-10-26).
  With `GOOGLE_SDK_GO_LOGGING_LEVEL=debug` the Google Cloud Go SDK
  logged unredacted token requests and responses and every KMS request
  and response, the keyset plaintext included (v1.0.0, v1.0.1). A
  discard logger (`option.WithLogger`) silences every SDK log line,
  present and future; a redacting logger would have to keep up with
  the SDK, and unsetting the variable would hide the intent.
- **Vault: token only in v1** (user, 2026-10-22). AppRole or Kubernetes
  auth is later work.
- **One credential per side** (user, 2026-10-22): `--credential` for the
  input, `--new-credential` for the output of `convert-keyset`, so it can
  move a keyset between backends.
- Encrypted I/O and encrypted conversion belong to this feature, so the
  KMS plan builds it whole (user, 2026-10-22).
- **AWS honours `AWS_PROFILE`** (user, 2026-10-23), with and without
  `--credential`, and only through the environment: no `--aws-profile`
  flag. The backend moved to `tink-go-awskms/v3` (aws-sdk-go-v2) for it:
  v2's aws-sdk-go v1 session read neither `~/.aws/config` nor a profile
  other than `default` from a credential file. Tink v3 still hard-codes
  `default` for an INI file, so circe loads INI files and the default
  chain itself and hands Tink the KMS client (`WithKMS`); CSV files keep
  Tink's handling. v2 and v3 both use the `associatedData` encryption
  context, so existing keysets still decrypt.
