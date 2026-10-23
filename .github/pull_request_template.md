<!--
The pull request title becomes the commit subject on main (squash merge),
and the release notes are built from those subjects. Start the title with
one of the types below, for example:

  fix: list-keyset prints disabled keys
  feat(kms): Vault AppRole auth

| Type         | Use it for                                      | Release notes      |
|--------------|-------------------------------------------------|--------------------|
| feat:        | A new command, flag or behaviour                | Features           |
| change:      | A change to existing behaviour                  | Changes            |
| fix:         | A bug fix                                       | Fixes              |
| hotfix:      | An urgent fix to the latest release             | Hotfixes           |
| emergency:   | An emergency patch (security, broken release)   | Emergency patches  |
| docs:        | Documentation only                              | left out           |
| test:        | Tests only                                      | left out           |
| chore:       | Build, CI, dependencies, tidying                | left out           |

A scope in brackets is optional: chore(deps):, fix(kms):.
-->

## What and why

<!-- What does this change, and why? Link the issue it closes: "Closes #123". -->

## Type

- [ ] `feat:` new feature
- [ ] `change:` change to existing behaviour
- [ ] `fix:` bug fix
- [ ] `hotfix:` urgent fix to the latest release
- [ ] `emergency:` emergency patch
- [ ] `docs:` documentation only (not in release notes)
- [ ] `test:` tests only (not in release notes)
- [ ] `chore:` build, CI, dependencies (not in release notes)

## Checklist

- [ ] The title starts with the type ticked above.
- [ ] `go test ./...`, `go vet ./...` and `go mod tidy -diff` pass.
- [ ] The matching feature in `features/` is updated, or a new one added.
- [ ] No keysets, key material, credentials or real key URIs anywhere in
      this pull request.
