# T09 · Pack detail page

| | |
|---|---|
| **Epic** | [E02 · Browse](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T08 |
| **Unlocks** | T10, T15 |
| **Spec** | §6.2 |
| **Status** | ☐ not started |

## Problem

`GET /packs/:slug` renders every **approved** sticker in a pack as a grid of tiles. This is
the page where copying happens, so the tile markup it emits is the contract T10 and T13
both build on.

## Given

- `pkg/ui/components/sticker.go` from T08
- Pack + sticker queries from T04
- `pkg/pager`

## Constraints

- Paginate. The global "Singles" pack is guaranteed to outgrow one page.
- Each tile must carry the data the client layer needs, as data attributes:
  `data-sticker-id`, `data-sticker-url`, `data-animated`, `data-name`. Get these right now
  — changing them later means touching three JS files.
- `data-animated` comes from the **database column**, never from the file extension.
- A 404 for an unknown slug, and for a pack with no approved stickers — not a blank page.
- No copy behaviour yet. Tiles render; clicking does nothing. That is T10.

## Acceptance

- [ ] `/packs/:slug` renders approved stickers only
- [ ] Pending and rejected stickers are absent from the HTML source, not merely hidden
- [ ] Unknown slug returns 404
- [ ] A pack whose stickers are all pending returns 404
- [ ] Every tile carries all four data attributes with correct values
- [ ] Pagination works

## Verify

```
make test && make build && make run
```

Then, against seeded data:

```
curl -s localhost:8000/packs/<slug> | grep -c 'data-sticker-id'
```

Compare with the approved count in `psql`. `view-source` and confirm no pending sticker's
object key appears anywhere in the page.

## Files

- `pkg/ui/pages/stickers.go`
- `pkg/ui/components/sticker.go`
- `pkg/handlers/stickers.go`
- `pkg/routenames/names.go`
- `pkg/handlers/stickers_test.go`

## Hints

- The "absent from source, not hidden" test is the one that catches a `LEFT JOIN` that
  should have been an `INNER JOIN`. Write it as a handler test.
- Keep `PackCard` and `StickerTile` in the same component file; they share styling.
