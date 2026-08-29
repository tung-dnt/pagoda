# T27 · Fix CI Go version + green pipeline

| | |
|---|---|
| **Epic** | [E08 · Ship](./README.md) |
| **Size** | XS — ~30 min |
| **Depends on** | None |
| **Unlocks** | T29 |
| **Spec** | Open Q1 |
| **Status** | ☐ not started |

## Problem

`.github/workflows/test.yml` pins `go-version: '1.25'` while `go.mod` declares `go 1.27.0`.
CI has been testing a different toolchain than the module targets. Fix it.

This is pre-existing, unrelated to the feature, and takes half an hour. Do it early — a
green pipeline you trust is worth more the longer the project runs.

## Given

- `.github/workflows/test.yml` — Go 1.25, `actions/checkout@v3`, `actions/setup-go@v3`,
  `actions/cache@v3`, Postgres 17 service

## Constraints

- Bump Go to match `go.mod` exactly.
- The v3 actions are several major versions behind; bumping them is in scope but keep it in
  a separate commit so a failure is attributable.
- Do not change `make test` itself here.
- CI must be green before you touch it further. If it is already red for an unrelated
  reason, fix that first or you cannot tell whether your change worked.

## Acceptance

- [ ] CI runs the same Go version as `go.mod`
- [ ] The workflow passes on a pushed branch
- [ ] The Postgres service still comes up healthy
- [ ] Dependency caching still hits

## Verify

Push a branch and watch the run:

```
gh run watch
```

## Files

- `.github/workflows/test.yml`

## Hints

- `go-version-file: go.mod` is better than a literal — it can never drift again. Prefer it.
