---
type: Feature
title: Contributing guide, issue forms and pull request template
description: CONTRIBUTING.md tells contributors how to report, propose and build changes; issue forms collect complete bug reports and feature requests; changes land as squash-merged feature branches whose typed PR title drives the release notes.
status: stable
code:
  - CONTRIBUTING.md
  - .github/ISSUE_TEMPLATE/**
  - .github/pull_request_template.md
tags: [project, github, contributing]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T23:00:00Z }
---

# Intent
circe is public open source. Contributors should be able to file a
useful bug report or feature request without being asked follow-up
questions, and to build, test and send a change without guessing.

# Behavior
- `CONTRIBUTING.md` covers, in order: what circe is and is not (not a
  Google tool; the Tink API is the contract, the Java tinkey CLI is not),
  never posting secrets, security reports, bug reports, feature
  requests, the development loop (`make`, `go test`, `go vet`,
  `golangci-lint run`, `go mod tidy`, `make dev` for AWS KMS on floci), how features in
  `features/` describe behaviour, pull requests, AI-assisted
  contributions, and the MIT licence of contributions.
- Security problems go to GitHub private vulnerability reporting
  (Security tab, "Report a vulnerability"), never to a public issue.
- New issues offer two forms and no blank issue:
  - **Bug report** (`bug_report.yml`, label `bug`): circe version
    (`circe --version`), install method, OS and architecture, KMS
    backend, the command run, expected and actual behaviour (stderr),
    steps to reproduce, and a required tick that no key material,
    credentials or real key URIs are in the report.
  - **Feature request** (`feature_request.yml`, label `enhancement`):
    the problem or use case, the proposed command or flags, whether Tink
    or tinkey already does it (with a link), alternatives, the area, and
    whether the reporter would send a pull request.
- `config.yml` turns blank issues off and adds two contact links: the
  private security report, and the README for usage questions.
- Every change is made on a feature branch and lands on `main` as one
  squash-merged commit whose subject is the pull request title; nothing
  is pushed to `main` directly.
- The pull request title starts with a type, with an optional scope in
  brackets (`fix(kms):`): `feat:`, `change:`, `fix:`, `hotfix:`,
  `emergency:` appear in the release notes (Features, Changes, Fixes,
  Hotfixes, Emergency patches); `docs:`, `test:`, `chore:` are left out
  ([Release builds](/build/release.md)). The pull request template
  (`.github/pull_request_template.md`) lists the types with a tick box
  each, and a checklist: title typed, tests, vet, lint and tidy pass, the
  feature updated, no key material.
- `CONTRIBUTING.md` has the same type table under Pull requests, and a
  section on Dependabot and govulncheck
  ([Dependency updates](/build/dependencies.md)).
- Labels used are GitHub's defaults (`bug`, `enhancement`), which the
  repository has.

# Interfaces
- [GitHub Actions CI](/build/ci.md) and [golangci-lint checks](/build/lint.md):
  the checks a pull request must pass.
- [Make targets](/build/makefile.md) and the
  [floci dev environment](/build/dev-floci.md): the development loop.
- [Command surface](/cli/command-surface.md): exit codes and errors a
  bug report quotes.

# Decisions
- **Contributing guide and issue templates** (user, 2026-10-23).
- **`CONTRIBUTING.md`, not `contribution.md`**: GitHub links that name
  from the new-issue and pull-request pages.
- **Issue forms (YAML), not markdown templates**: required fields and
  dropdowns, so reports arrive complete.
- **A required no-secrets tick**: a keyset pasted into a public issue
  leaks its key; `list-keyset` output has no key material, so the forms
  ask for that instead.
- **Feature branches, squash merged into `main`** (user, 2026-10-23):
  one commit per change on `main`, so the commit subjects are the
  release notes.
- **Types `change`, `hotfix` and `emergency`** (user, 2026-10-23) beside
  `feat`, `fix` and the excluded `docs`, `test`, `chore`. `emergency:` is
  the prefix for an emergency patch.
- **AI-assisted contributions follow the README's stance**: welcome,
  with the contributor responsible for every line.
