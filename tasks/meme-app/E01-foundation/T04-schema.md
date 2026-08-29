# T04 · Schema migration + sqlc queries

| | |
|---|---|
| **Epic** | [E01 · Foundation](./README.md) |
| **Size** | M — ~3 h |
| **Depends on** | None |
| **Unlocks** | T07, T08, T17 |
| **Spec** | §6.1 |
| **Status** | ☐ not started |

## Problem

Create the `packs`, `stickers`, `tags`, and `pack_tags` tables, seed the one global
"Singles" pack, and write the sqlc queries the rest of the app runs on.

The single most important property: **a visitor can never see a non-approved sticker**, and
that is enforced in SQL, not in handler code.

## Given

- `pkg/postgres/migrations/000001_init.up.sql` — the pattern to follow
- `pkg/postgres/queries/users.sql` — the sqlc query comment style
- `sqlc.yaml` — reads migrations as schema, emits to `pkg/postgres/db`
- Full DDL in [spec §6.1](../../../docs/spec/phase-1.md)

## Constraints

- **At most one global pack**, enforced by the database:
  `CREATE UNIQUE INDEX packs_one_global_idx ON packs (is_global) WHERE is_global;`
  A partial unique index is how Postgres expresses "at most one row where X".
- `status` lives on **stickers**, never on packs. Pack visibility is derived from a join
  and never stored, so the two cannot drift.
- Approved-sticker counts are **computed** (`COUNT(*) FILTER (WHERE status = 'approved')`),
  not denormalized. A counter column would need maintaining on approve, reject, and delete.
- Write a reversible `.down.sql`. Test it.
- `pkg/postgres/db/` is generated — never hand-edit it.

## Acceptance

- [ ] `make migrate-up` then `make migrate-down` then `make migrate-up` all succeed
- [ ] Inserting a second row with `is_global = TRUE` fails
- [ ] `make sqlc-gen` produces compiling code; `make sqlc-vet` is clean
- [ ] Queries exist for: list visible packs (paginated), get pack by slug, list a pack's
      approved stickers (paginated), get sticker by id, count approved per pack
- [ ] Every visitor-facing query filters `stickers.status = 'approved'` in the SQL itself

## Verify

```
make db-up
make migrate-up && make migrate-down && make migrate-up
make sqlc-gen && make sqlc-vet
make test
```

Plus a query test asserting a pending sticker is absent from every visitor query.

## Files

- `pkg/postgres/migrations/000002_stickers.up.sql` / `.down.sql`
- `pkg/postgres/queries/packs.sql`, `stickers.sql`, `tags.sql`
- `pkg/postgres/db/*` (generated)
- `pkg/postgres/db/..._test.go`

## Hints

- `make migrate-new name=stickers` creates the pair with the right numbering.
- sqlc names the Go method from the `-- name: X :many` comment; keep names verb-first.
- A pack with zero approved stickers must not appear — that is an `INNER JOIN` or an
  `EXISTS`, not a `LEFT JOIN`. Write the test that proves it.
