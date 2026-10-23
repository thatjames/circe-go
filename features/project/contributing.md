---
type: Feature
title: Contributing guide and issue forms
description: CONTRIBUTING.md tells contributors how to report, propose and build changes; GitHub issue forms collect bug reports and feature requests with the details triage needs and no key material.
status: stable
code:
  - CONTRIBUTING.md
  - .github/ISSUE_TEMPLATE/**
tags: [project, github, contributing]
generated: { by: claude-code/claude-opus-5-5, at: 2026-10-23T21:15:24Z }
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
  `go mod tidy`, `make dev` for AWS KMS on floci), how features in
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
- Labels used are GitHub's defaults (`bug`, `enhancement`), which the
  repository has.

# Interfaces
- [GitHub Actions CI](/build/ci.md): the checks a pull request must pass.
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
- **AI-assisted contributions follow the README's stance**: welcome,
  with the contributor responsible for every line.
