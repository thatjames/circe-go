# circe

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

circe is pure Go (built with `CGO_ENABLED=0`) and needs Go 1.26 or later.
It supports Linux and macOS.

## Usage

```sh
circe list-key-templates                                  # Tink template names
circe create-keyset --key-template AES256_GCM --out ks.json
circe list-keyset --in ks.json                            # key IDs, types, statuses; no key material
circe convert-keyset --in ks.json --out ks.bin --out-format binary
```

Keysets are read from `--in` (default stdin) and written to `--out`
(default stdout) in Tink's `json` or `binary` format (`--in-format`,
`--out-format`). `--out` never overwrites an existing file. Errors exit
1; usage errors exit 2.

More commands land as circe grows.

## Licence

MIT, see [LICENSE](LICENSE).
