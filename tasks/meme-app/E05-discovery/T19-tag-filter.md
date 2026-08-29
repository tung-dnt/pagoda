# T19 · Tag filtering

| | |
|---|---|
| **Epic** | [E05 · Discovery](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T18 |
| **Unlocks** | — |
| **Spec** | §6.2 |
| **Status** | ☐ not started |

## Problem

Make tag chips clickable: `GET /search?tag=slug` filters to packs carrying that tag, and
combines with `q` when both are present.

## Given

- Search from T18, chips from T17

## Constraints

- `q` and `tag` compose — `?q=cat&tag=reaction` means both, not either.
- The active tag must be visibly active, and there must be an obvious way to clear it.
- Filter state lives in the URL, so a filtered view is shareable and the back button works.
- Unknown tag slug → empty state, not a 404 and not an error.
- Same approved-only rule as everything else.

## Acceptance

- [ ] Clicking a chip filters results
- [ ] `q` + `tag` together apply both
- [ ] The active tag is visually distinct and clearable
- [ ] Back button restores the previous filter
- [ ] Unknown tag renders an empty state
- [ ] Pagination preserves both params

## Verify

```
make run
```

Click a chip → filtered. Add a search term → both apply. Back → previous state. Copy the
URL into a new tab → identical view.

## Files

- `pkg/handlers/search.go`
- `pkg/ui/pages/search.go`
- `pkg/postgres/queries/packs.sql`

## Hints

- Build the query string in one helper used by chips *and* pagination, or the two will
  disagree about which params to keep.
