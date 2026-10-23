# circe

[![CI](https://github.com/thatjames/circe-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thatjames/circe-go/actions/workflows/ci.yml)
[![Vulnerabilities](https://github.com/thatjames/circe-go/actions/workflows/vuln.yml/badge.svg?branch=main)](https://github.com/thatjames/circe-go/actions/workflows/vuln.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/thatjames/circe-go.svg)](https://pkg.go.dev/github.com/thatjames/circe-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/thatjames/circe-go)](https://goreportcard.com/report/github.com/thatjames/circe-go)
[![Release](https://img.shields.io/github/v/release/thatjames/circe-go)](https://github.com/thatjames/circe-go/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/thatjames/circe-go)](go.mod)
[![License: MIT](https://img.shields.io/github/license/thatjames/circe-go)](LICENSE)

## Important Notice

This is **NOT** an official Google tool, it is only based off of Google's APIs. I am just some dude, use your common sense when using random internet software.

## AI Disclosure

Every line of code (for all intents and purposes) was written by claude. I am a senior engineer and I am fully responsible for the code in this repository. Claude may have typed it, but I made the decisions. 

## What is circe?

`circe` is a command-line tool for [Tink](https://developers.google.com/tink)
keysets, written in Go on [tink-go](https://github.com/tink-crypto/tink-go).
It does what Tink's `tinkey` does: create and manage keysets, inspect and
convert them, and keep them encrypted with a key in GCP KMS, AWS KMS or
HashiCorp Vault.

## Tink's API, not Java tinkey's CLI

circe follows the Tink API, Tink's JSON and binary keyset formats and
Tink's key template names. A keyset written by circe can be read by any
Tink library, and the other way round. circe's commands and flags may
differ from those of the Java `tinkey` CLI. Where there is no reason to
differ, they keep tinkey's names.

## Install

```sh
go install github.com/thatjames/circe-go/cmd/circe@latest
```

circe is pure Go (built with `CGO_ENABLED=0`) and needs Go 1.27 or later.
It supports Linux and macOS.

## Usage

```sh
circe list-key-templates                                  # Tink template names
circe create-keyset --key-template AES256_GCM --out ks.json
circe list-keyset --in ks.json                            # key IDs, types, statuses; no key material
circe convert-keyset --in ks.json --out ks.bin --out-format binary

# Rotate: add a key, then make it the primary.
circe add-key --key-template AES256_GCM --in ks.json --out ks2.json
circe promote-key --key-id 1234567 --in ks2.json --out ks3.json
circe disable-key --key-id 7654321 --in ks3.json --out ks4.json   # also enable-key, delete-key

# The public half of a signature, hybrid or JWT signature keyset.
circe create-keyset --key-template ED25519 --out priv.json
circe create-public-keyset --in priv.json --out pub.json
```

### Encrypted keysets

Every command takes `--master-key-uri` (and `--credential`, a credentials
file) to read and write keysets encrypted with a KMS key, as Tink does:

```sh
circe create-keyset --key-template AES256_GCM \
  --master-key-uri gcp-kms://projects/p/locations/global/keyRings/r/cryptoKeys/k --out ks.json
circe list-keyset --in ks.json --master-key-uri gcp-kms://projects/p/locations/global/keyRings/r/cryptoKeys/k

# Move a keyset to another KMS key, or out of KMS (omit --new-master-key-uri).
circe convert-keyset --in ks.json --master-key-uri gcp-kms://... \
  --new-master-key-uri aws-kms://arn:aws:kms:eu-west-1:111122223333:key/... --out ks-aws.json
```

| Backend | URI | Credentials without `--credential` |
|---------|-----|------------------------------------|
| GCP KMS | `gcp-kms://projects/.../cryptoKeys/...` | Application Default Credentials |
| AWS KMS | `aws-kms://arn:aws:kms:...` | the AWS SDK default chain, with the profile `AWS_PROFILE` names |
| HashiCorp Vault | `hcvault://host:8200/transit/keys/<name>` | `VAULT_TOKEN` (token auth only) |

AWS uses the profile named by `AWS_PROFILE` (else `default`) from
`~/.aws/config` and `~/.aws/credentials`, SSO and assume-role profiles
included, as the AWS CLI does. The region comes from the key URI.
`--credential` takes either a CSV file of access keys as the IAM console
downloads it, or a shared credentials/config file. From a shared file
circe uses the `AWS_PROFILE` profile (else `default`), and a profile
missing from that file is an error:

```sh
AWS_PROFILE=prod circe list-keyset --in ks.json --master-key-uri aws-kms://arn:aws:kms:eu-west-1:111122223333:key/...
```

`convert-keyset` reads both `--credential` and `--new-credential` with
the same `AWS_PROFILE`. To move a keyset between two AWS profiles, pass
each side its own file.

### Files and errors

Keysets are read from `--in` (default stdin) and written to `--out`
(default stdout) in Tink's `json` or `binary` format (`--in-format`,
`--out-format`). `--out` never overwrites an existing file. Errors exit
1; usage errors exit 2.

More commands land as circe grows.

## Licence

MIT, see [LICENSE](LICENSE).
