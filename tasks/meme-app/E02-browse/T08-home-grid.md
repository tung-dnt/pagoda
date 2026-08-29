# T08 · Home page: pack grid

| | |
|---|---|
| **Epic** | [E02 · Browse](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T04, T07 |
| **Unlocks** | T09, T17 |
| **Spec** | §6.2 |
| **Status** | ☐ not started |

## Problem

Replace the boilerplate home page with a responsive grid of pack cards. Each card shows a
cover sticker, the pack name, and its approved-sticker count. A pack appears only if it has
at least one approved sticker.

This is the first task where the app looks like the product.

## Given

- `pkg/ui/pages/home.go` — currently renders demo posts; rewrite it
- `pkg/ui/layouts/primary.go` — the page shell
- `pkg/ui/components/` — existing DaisyUI-flavoured building blocks
- `pkg/pager` — `NewPager(ctx, itemsPerPage)`, already wired for the `page` query param
- Queries from T04

## Constraints

- Follow the house view style exactly: dot-imported gomponents, functions returning `Node`,
  params passed as a struct. See `pkg/ui/forms/file.go`.
- Route names go in `pkg/routenames` — never inline a path in a view.
- Pack visibility comes from the SQL, not a filter in the handler.
- The global "Singles" pack renders as an ordinary card.
- Images are lazy-loaded (`loading="lazy"`) and must not cause layout shift — set explicit
  `width`/`height` from the stored dimensions.
- Cards must be tappable targets on mobile, not just hover-able on desktop.

## Acceptance

- [ ] `/` renders a grid of pack cards from real seeded data
- [ ] A pack with zero approved stickers does not appear
- [ ] Pending and rejected stickers are never used as a cover image
- [ ] Counts match `COUNT(*) FILTER (WHERE status = 'approved')`
- [ ] Pagination works and preserves the `page` query param
- [ ] The grid is usable at 375 px wide

## Verify

```
make test && make build && make run
```

Open `/`, compare against `psql` counts. Resize to 375 px. Check DevTools for CLS.

## Files

- `pkg/ui/pages/home.go`
- `pkg/ui/components/sticker.go` (new — `PackCard`)
- `pkg/handlers/stickers.go` (new)
- `pkg/routenames/names.go`
- `pkg/handlers/stickers_test.go`

## Hints

- Handlers self-register: `func init() { Register(new(Stickers)) }`, then `Init` and
  `Routes`. Copy the shape from `pkg/handlers/files.go`.
- Cover sticker = lowest `position` among approved. Do it in SQL with `DISTINCT ON`.
- The handler test should assert a pending-only pack is absent from the response body.
