# T15 · Pack-level favorite

| | |
|---|---|
| **Epic** | [E04 · Favorites](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T14 |
| **Unlocks** | — |
| **Spec** | §6.4 |
| **Status** | ☐ not started |

## Problem

Add a "favorite this pack" control on the pack detail page. It stores the pack plus every
approved sticker's blob — **except on the global "Singles" pack**, where the control is not
rendered at all.

## Given

- `db.js` (`favorite_packs` store), `favorites.js`, the Alpine `favorites` store from T13
- Pack detail page from T09

## Constraints

- **The global pack is excluded.** It grows without bound, so one tap would try to cache
  every loose sticker ever uploaded. Spec §8 lists rendering this control on the global
  pack under **Never**. Gate on `is_global` server-side — do not hide it with CSS.
- Caching a whole pack is N fetches. Show progress, and make it cancellable. A user on
  mobile data deserves to know what they started.
- Progress is Alpine component state (`done`, `total`, `cancelling`) bound with `x-text`
  and `x-show`. An `AbortController` held in that state handles cancel. Do not drive a
  progress bar with manual DOM writes.
- Respect the 100 MB cap from T12. If a pack would exceed it, say so and stop rather than
  silently evicting the user's older favorites to make room.
- Favoriting a pack should mark its individual stickers as favorited too — one concept, not
  two competing ones.
- Unfavoriting a pack removes the pack record. Decide and document what happens to
  individually-favorited stickers inside it; do not leave it ambiguous.

## Acceptance

- [ ] The pack page shows a favorite-pack control
- [ ] The "Singles" pack renders no such control, and none appears in the HTML source
- [ ] Favoriting a pack caches all its approved stickers with visible progress
- [ ] Cancelling mid-way leaves a consistent state, not half a pack
- [ ] Exceeding the cap shows a clear message and does not evict existing favorites
- [ ] Favorited packs appear on `/favorites`

## Verify

```
make run
```

Favorite a 10-sticker pack, watch progress, go offline, confirm all 10 render on
`/favorites`. Then open `/packs/singles` and `view-source` — grep for the control, expect
nothing. Then temporarily drop the cap to 1 MB and confirm the refusal message.

## Files

- `public/static/js/favorites.js`
- `pkg/ui/components/sticker.go`
- `pkg/ui/pages/stickers.go`

## Hints

- Sequential fetches with a progress counter beat `Promise.all` here — you get progress and
  cancellation for free, and you will not hammer the bucket with 200 parallel requests.
- An `AbortController` handles cancel cleanly.
