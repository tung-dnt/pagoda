# T07 · Dev seed command

| | |
|---|---|
| **Epic** | [E02 · Browse](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T04, T05, T06 |
| **Unlocks** | T08 |
| **Spec** | — |
| **Status** | ☐ not started |

## Problem

You cannot build a browse UI against an empty database. Add a dev-only command that
creates a handful of packs with real images — including animated ones and a few pending
stickers, so every visibility rule is observable by eye.

## Given

- `cmd/admin/main.go` — the pattern for a small CLI that boots the container
- `services.Storage` from T05, `Inspect` from T06

## Constraints

- Dev-only. It must refuse to run when `app.environment` is not `local`.
- Seed data must include, at minimum:
  - 3 normal packs with 5–10 stickers each
  - the global "Singles" pack with several loose stickers
  - **at least one animated sticker per pack** — otherwise T10's second copy path is
    untested by eye
  - **at least one pending and one rejected sticker** — so you can confirm they are absent
    from browse
  - one pack where *every* sticker is pending, which must not appear in the grid at all
- Idempotent: running it twice must not duplicate. Truncate-and-reseed is fine in dev.
- Images go through the same `Put` → `Promote` path real uploads use. Do not write rows
  pointing at objects that were never uploaded.

## Acceptance

- [ ] `make seed` populates the database and the object store
- [ ] Running it twice leaves the same state as running it once
- [ ] It exits with an error when `environment != local`
- [ ] The seeded set covers every case in Constraints above

## Verify

```
make db-up && make migrate-up && make seed
psql "$DATABASE_URL" -c "SELECT status, count(*) FROM stickers GROUP BY status;"
```

Expect all three statuses present.

## Files

- `cmd/seed/main.go`
- `Makefile` (add the `seed` target)
- `cmd/seed/testdata/*` or a small set of committed sample images

## Hints

- Pull a few genuinely funny stickers you actually want in the app — you will be looking
  at this data for the next five weeks.
- Keep the fixture images small; they are committed.
