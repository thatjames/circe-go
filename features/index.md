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

* [Keyset I/O](/keyset/keyset-io.md) - Reads and writes Tink keysets in Tink's JSON and binary formats, cleartext or KMS-encrypted, through tink-go only.
* [Key templates](/keyset/key-templates.md) - Maps Tink key template names (AES256_GCM, ECDSA_P256, ...) to tink-go templates for --key-template, and lists them.
* [Keyset lifecycle](/keyset/lifecycle.md) - create-keyset, add-key, promote-key, rotate-keyset, enable-key, disable-key, destroy-key and delete-key manage keys in a Tink keyset.
* [Inspect and convert keysets](/keyset/inspect-convert.md) - list-keyset shows a keyset's key info without key material, convert-keyset changes format or encryption, create-public-keyset extracts the public keyset.

# KMS

* [KMS-encrypted keysets](/kms/encrypted-keysets.md) - --master-key-uri encrypts and decrypts keysets with a key in GCP KMS, AWS KMS or HashiCorp Vault, via the tink-go KMS extensions.
