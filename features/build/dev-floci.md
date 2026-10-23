---
type: Feature
title: floci dev environment
description: compose.yaml runs floci, a local AWS emulator on localhost:4566, so AWS KMS keysets can be tried without an AWS account; make dev starts it.
status: stable
code:
  - compose.yaml
tags: [build, dev, aws, floci]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T20:57:54Z }
---

# Intent
Try circe's `aws-kms://` keysets locally against an emulated KMS, with
no AWS account and no real credentials.

# Behavior
- `compose.yaml` (compose project `circe`) has one service, `floci`:
  image `floci/floci` pinned to a version, container `circe-floci`,
  port `4566` published on `127.0.0.1` only.
- floci runs with `FLOCI_STORAGE_MODE=persistent` (its default,
  `memory`, loses every key on stop) and keeps its state in the named
  volume `floci-data` (mounted at `/app/data`), so KMS keys created in it survive `make dev-down` and a
  restart; `docker compose down -v` wipes it.
- The image's own healthcheck decides readiness; `make dev` waits for it.
- Point circe at it through the aws-sdk-go-v2 environment, for example
  `AWS_ENDPOINT_URL=http://localhost:4566 AWS_REGION=us-east-1
  AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test`, and a key URI
  `aws-kms://arn:aws:kms:us-east-1:000000000000:key/<id>` from
  `aws kms create-key` against the same endpoint.
- CI does not start it: tests stay offline ([CI](/build/ci.md)).

# Interfaces
- [Make targets](/build/makefile.md): `make dev`, `make dev-down`.
- [KMS-encrypted keysets](/kms/encrypted-keysets.md): the AWS backend
  loads the default config chain, which reads `AWS_ENDPOINT_URL`.

# Decisions
- **floci, from compose** (user, 2026-10-23).
- Bound to `127.0.0.1`: the emulator accepts any credentials.
- **`persistent`, not `hybrid`**: every write reaches disk before the
  response, so a key created just before `make dev-down` is not lost.
- No Docker socket mount: circe needs KMS only, not floci's container-
  backed services (Lambda, RDS).
