---
okf_version: "0.2"
---

# Build

* [Go module and toolchain](/build/go-module.md) - circe-go is the pure-Go module github.com/thatjames/circe-go on Go 1.26, depending on tink-go and cobra only through go.mod.
* [GitHub Actions CI](/build/ci.md) - On every push and pull request, GitHub Actions runs go vet and go test with cgo off on Linux and macOS.
* [Release builds](/build/release.md) - Pushing a vX.Y.Z tag builds static circe binaries for Linux and macOS and publishes them as a GitHub release.

# CLI

* [Command surface](/cli/command-surface.md) - The circe binary is a cobra CLI with one subcommand per tinkey operation, shared keyset input/output flags, and errors on stderr with exit 1.

# Keyset

* [Keyset I/O](/keyset/keyset-io.md) - Reads and writes cleartext Tink keysets in Tink's JSON and binary formats through tink-go only; the hook for KMS-encrypted keysets.
* [Key templates](/keyset/key-templates.md) - Maps Tink key template names (AES256_GCM, ECDSA_P256, ...) to tink-go templates for --key-template, and lists them.
* [Create keyset](/keyset/create-keyset.md) - create-keyset writes a new Tink keyset with one enabled primary key from a named Tink key template.
* [Key management](/keyset/key-management.md) - add-key, promote-key, enable-key, disable-key and delete-key change the keys of an existing Tink keyset; no destroy-key and no rotate-keyset.
* [Create public keyset](/keyset/public-keyset.md) - create-public-keyset reads a private signature, hybrid or JWT signature keyset and writes its public keyset in cleartext.
* [Inspect and convert keysets](/keyset/inspect-convert.md) - list-keyset shows a keyset's key info without key material; convert-keyset rewrites a keyset between Tink's JSON and binary formats.

# KMS

* [KMS-encrypted keysets](/kms/encrypted-keysets.md) - --master-key-uri reads and writes KMS-encrypted keysets on every command (GCP KMS, AWS KMS, Vault), and convert-keyset moves keysets between cleartext and KMS keys.
