---
okf_version: "0.2"
---

# Build

* [Go module and toolchain](/build/go-module.md) - circe-go is the pure-Go module github.com/thatjames/circe-go on Go 1.27, depending on tink-go and cobra only through go.mod.
* [GitHub Actions CI](/build/ci.md) - On every push and pull request, GitHub Actions runs go vet, go test and make with cgo off on Linux and macOS.
* [Release builds](/build/release.md) - Pushing a vX.Y.Z tag builds static circe binaries for Linux and macOS, publishes them as a GitHub release, and pushes the Docker image to ghcr.io/thatjames/circe.
* [Make targets](/build/makefile.md) - A Makefile builds the circe binary (default), the Docker image tagged with the latest git tag (docker, docker-push) and starts the floci dev container (dev); CI calls it.
* [Docker image](/build/docker-image.md) - A multi-stage Dockerfile builds a static circe binary for linux/amd64 and linux/arm64 and ships it alone on distroless static (non-root, CA certificates), entrypoint circe.
* [floci dev environment](/build/dev-floci.md) - compose.yaml runs floci, a local AWS emulator on localhost:4566, so AWS KMS keysets can be tried without an AWS account; make dev starts it.
* [Dependency updates and vulnerability checks](/build/dependencies.md) - Dependabot opens weekly grouped update pull requests for Go modules, GitHub Actions and Docker base images; govulncheck fails a Vulnerabilities workflow when circe calls a known-vulnerable symbol.
* [golangci-lint checks](/build/lint.md) - A Lint workflow runs golangci-lint v2 with a curated .golangci.yml on every push and pull request and fails on any finding; its badge replaces the retired Go Report Card badge on the README.

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

# Primitives

* [Primitive commands](/primitives/primitive-commands.md) - Shared shape of circe's per-primitive commands (aead, daead, streaming-aead, mac, prf, signature, hybrid, JWT) - --keyset and its KMS flags, data on --in/--out, --encoding raw|base64|hex, and wrong-keyset and verification errors.
* [AEAD encrypt and decrypt](/primitives/aead.md) - circe aead encrypt and decrypt run Tink's AEAD primitive on --in with an AEAD keyset and optional associated data.
* [Deterministic AEAD encrypt and decrypt](/primitives/deterministic-aead.md) - circe daead encrypt and decrypt run Tink's deterministic AEAD on --in with a DAEAD keyset; equal input and associated data give an equal ciphertext.
* [Streaming AEAD encrypt and decrypt](/primitives/streaming-aead.md) - circe streaming-aead encrypt and decrypt stream --in to --out in constant memory with Tink's streaming AEAD; raw bytes only, and a failed decrypt leaves no --out file.
* [MAC compute and verify](/primitives/mac.md) - circe mac compute writes the Tink MAC tag of --in; circe mac verify checks a --tag against --in and exits 0, or 3 when it does not match.
* [PRF compute](/primitives/prf.md) - circe prf compute --output-length <n> writes n bytes of the keyset's primary PRF over --in.
* [Signature sign and verify](/primitives/signature.md) - circe signature sign signs --in with a private signature keyset; circe signature verify checks a --signature against --in with the public keyset and exits 0, or 3 when it does not verify.
* [Hybrid encrypt and decrypt](/primitives/hybrid.md) - circe hybrid encrypt encrypts --in to a public hybrid keyset with optional context info; circe hybrid decrypt decrypts with the private keyset.
* [JWT MAC and signature](/primitives/jwt.md) - circe signs or MACs a JWT claims JSON into a compact JWT with a Tink JWT keyset, and verifies a compact JWT against Tink validator flags, printing its verified claims.

# Project

* [Contributing guide, issue forms and pull request template](/project/contributing.md) - CONTRIBUTING.md tells contributors how to report, propose and build changes; issue forms collect complete bug reports and feature requests; changes land as squash-merged feature branches whose typed PR title drives the release notes.
