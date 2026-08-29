# T18 · Search endpoint + results

| | |
|---|---|
| **Epic** | [E05 · Discovery](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T17 |
| **Unlocks** | T19 |
| **Spec** | §6.2 |
| **Status** | ☐ not started |

## Problem

Replace the boilerplate's fake search with a real one: `GET /search?q=` matching pack
names, sticker names, and tag names — over approved content only.

## Given

- `pkg/handlers/search.go` — currently returns lorem ipsum; replace it
- `pkg/ui/pages/search.go`
- `pkg/pager`

## Constraints

- `ILIKE '%' || $1 || '%'` is deliberately chosen over full-text search. `pg_trgm` is
  **deferred** until a measured query is actually slow — see spec Open Question 2.
- Escape `%` and `_` in user input, or a query containing `%` matches everything.
- Restricted to `stickers.status = 'approved'`, enforced in the SQL.
- Results mix packs and stickers. Decide the ranking now and write it down — name-prefix
  matches above substring matches is a reasonable, explainable rule.
- Empty query renders the browse page or an empty state, never an error.
- Paginate.

## Acceptance

- [ ] `/search?q=cat` returns packs, stickers, and tag-matched packs containing "cat"
- [ ] No pending or rejected sticker ever appears in results
- [ ] `q=%` returns nothing pathological
- [ ] Empty `q` renders cleanly
- [ ] Results are paginated and ranked by the documented rule
- [ ] The lorem ipsum handler is gone

## Verify

```
make test && make run
```

Search a term you know is only on a pending sticker → zero results. Search `%` → sane.
Search a term with 30+ matches → pagination works.

## Files

- `pkg/handlers/search.go`, `pkg/handlers/search_test.go`
- `pkg/ui/pages/search.go`
- `pkg/postgres/queries/packs.sql`, `stickers.sql`

## Hints

- Write the "term only on a pending sticker returns zero results" test first. It is the
  one that matters and the easiest to forget.
- HTMX makes live-as-you-type search a two-line change — `hx-trigger="keyup changed delay:300ms"`
  — but only add it once the plain form works.
