# T17 · Tag queries + chips

| | |
|---|---|
| **Epic** | [E05 · Discovery](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T04, T08 |
| **Unlocks** | T18, T19 |
| **Spec** | §6.2 |
| **Status** | ☐ not started |

## Problem

Surface tags: queries to list them with counts, a chip component, and tag display on pack
cards and the pack page. Filtering comes in T19; this task just makes tags visible.

## Given

- `tags` / `pack_tags` tables from T04
- `pkg/ui/components/` for the chip

## Constraints

- Tag counts must reflect **visible** packs only — a tag attached solely to packs with no
  approved stickers should not appear, or should not claim a count it cannot deliver.
- Tags are attached to packs, not to individual stickers. Keep it that way in phase 1.
- Slugs are normalized on write: lowercase, trimmed, spaces to hyphens. Do it in one place
  so upload (T21) and any admin editing agree.
- Do not build the filter UI here. One task, one concern.

## Acceptance

- [ ] Query: list all tags with their visible-pack counts, ordered by count
- [ ] Query: list a pack's tags
- [ ] `TagChip` component renders consistently on cards and the pack page
- [ ] A tag with no visible packs does not appear in the list
- [ ] `make sqlc-vet` clean

## Verify

```
make sqlc-gen && make sqlc-vet && make test && make run
```

Seed a tag on a pack whose stickers are all pending; confirm it is absent from the tag list.

## Files

- `pkg/postgres/queries/tags.sql`
- `pkg/ui/components/sticker.go`
- `pkg/ui/pages/home.go`, `stickers.go`
- `pkg/postgres/db/..._test.go`

## Hints

- The count query needs the same "pack has ≥1 approved sticker" join as T04's pack list.
  Consider a SQL view so the rule lives in exactly one place.
