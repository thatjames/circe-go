# Contributing to circe

Thanks for helping. Bug reports, feature requests and pull requests are
all welcome. This page says how to make each one easy to act on.

## Before you start

- circe is **not** an official Google tool. Problems with Tink itself
  belong in [tink-go](https://github.com/tink-crypto/tink-go/issues).
- circe follows the **Tink API**: its keyset formats and key template
  names. A keyset circe writes must be readable by every Tink library.
  circe's commands and flags may differ from the Java `tinkey` CLI, so
  "tinkey does it differently" is not a bug by itself; "Tink can't read
  this keyset" is.

## Never post secrets

Issues and pull requests are public. Never paste:

- a keyset (cleartext or encrypted), key material or `keyData` values;
- credentials, tokens, or the contents of a `--credential` file;
- real KMS key URIs, account IDs or project names.

Use `circe list-keyset` output instead (it shows key IDs, types and
statuses, never key material), or a keyset you made just for the report.
Replace real URIs with placeholders such as
`aws-kms://arn:aws:kms:eu-west-1:111122223333:key/EXAMPLE`.

## Security problems

Report them privately through
[GitHub's private vulnerability reporting](https://github.com/thatjames/circe-go/security/advisories/new)
(the repository's Security tab, "Report a vulnerability"). Please don't
open a public issue for them.

## Reporting a bug

Open a [bug report](https://github.com/thatjames/circe-go/issues/new?template=bug_report.yml).
The form asks for:

- `circe --version` and how you installed circe;
- your OS and architecture;
- the KMS backend, if any;
- the exact command, what you expected, and what happened (the stderr
  output and the exit code);
- the smallest steps that reproduce it, ideally from
  `circe create-keyset` onwards.

## Proposing a feature

Open a [feature request](https://github.com/thatjames/circe-go/issues/new?template=feature_request.yml).
Start with the problem, not the solution: what are you trying to do, and
what stops you? Say whether Tink or `tinkey` already does it, and link
to it if so. For anything larger than a small fix, please open the
request before writing code, so we can agree on the shape first.

## Development

You need Go 1.27 or later, `make`, and Docker for the image and the dev
environment.

```sh
make                 # build ./circe (static, CGO_ENABLED=0)
go test ./...
go vet ./...
golangci-lint run    # lint and format, config in .golangci.yml
go mod tidy -diff    # must print nothing
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
make docker          # build the image circe:<latest tag>
```

CI runs the same checks on Linux and macOS for every push and pull
request (golangci-lint on Linux only); all of them must pass. Install
the golangci-lint version pinned in `.github/workflows/lint.yml`
(<https://golangci-lint.run/docs/welcome/install/>). Fix a finding
rather than silence it; when a rule is wrong for one line, add
`//nolint:<linter> // <reason>` there.

Tests must not need a cloud account or credentials. Use synthetic
keysets created in the test and fake AEADs for KMS.

### Trying AWS KMS locally

`make dev` starts [floci](https://floci.io), a local AWS emulator, on
`localhost:4566`. Point circe at it with throwaway credentials:

```sh
make dev
export AWS_ENDPOINT_URL=http://localhost:4566 AWS_REGION=us-east-1 \
       AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test
unset AWS_PROFILE
ARN=$(aws kms create-key --query KeyMetadata.Arn --output text)
./circe create-keyset --key-template AES256_GCM --master-key-uri "aws-kms://$ARN" --out ks.json
./circe list-keyset --in ks.json --master-key-uri "aws-kms://$ARN"
make dev-down        # stops floci; its keys are kept for next time
```

## Features describe the behaviour

Every part of circe is described by a feature document in
[`features/`](features/index.md): what it does, its edge cases and
errors, and the decisions behind it. A pull request that changes
behaviour updates the matching feature in the same change, or adds a new
one. If you're unsure which feature that is, ask in the pull request.

## Pull requests

Every change is made on a feature branch and lands on `main` as one
squash-merged commit; nothing is pushed to `main` directly.

- Keep each pull request to one change, with tests.
- The pull request title becomes the commit subject on `main`, and the
  release notes are built from those subjects. Start it with a type:

  | Type | Use it for | Release notes |
  |------|------------|---------------|
  | `feat:` | a new command, flag or behaviour | Features |
  | `change:` | a change to existing behaviour | Changes |
  | `fix:` | a bug fix | Fixes |
  | `hotfix:` | an urgent fix to the latest release | Hotfixes |
  | `emergency:` | an emergency patch (security, broken release) | Emergency patches |
  | `docs:` | documentation only | left out |
  | `test:` | tests only | left out |
  | `chore:` | build, CI, dependencies, tidying | left out |

  A scope is optional: `fix(kms): ...`, `chore(deps): ...`. The pull
  request template lists the types.
- Use the description for the why, and link the issue it closes.
- Releases are tagged by the maintainer; don't change version numbers.

## Dependencies and vulnerabilities

Dependabot opens weekly pull requests for Go modules (`fix(deps):`, in
the release notes), GitHub Actions and the Docker base images
(`chore(deps):`, left out). `govulncheck` runs on every push and pull
request and every Monday, so a newly published vulnerability in a
dependency fails the Vulnerabilities check on `main`. Run it locally
with:

```sh
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## AI-assisted contributions

Using an AI assistant is fine; circe itself is written that way (see the
README). You are responsible for every line you send: understand it,
test it, and be ready to explain it.

## Licence

circe is MIT licensed. By contributing you agree that your contribution
is licensed under the [MIT licence](LICENSE).
