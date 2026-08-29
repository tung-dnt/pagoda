# T13 · favorites.js: toggle + blob cache

| | |
|---|---|
| **Epic** | [E04 · Favorites](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T12, T11 |
| **Unlocks** | T14 |
| **Spec** | §6.4 |
| **Status** | ☐ not started |

## Problem

Add a favorite button to every sticker tile. Toggling it stores the sticker's metadata
**and fetches and caches its image blob**, so the favorites page later renders with no
network round-trip.

## Given

- `db.js` from T12
- The Alpine `stickerGrid` component and delegation pattern from T11

## Constraints

- Favoriting is **metadata + blob**. Metadata alone would make T14's offline requirement
  impossible.
- Fetch the blob with `mode: 'cors'`; reuse whatever T10 already proved works.
- Favorite state belongs in a single **`Alpine.store('favorites')`**, not in per-component
  state. The grid, the pack page, and the favorites page all read the same store, so a
  toggle in one place updates every other view for free. Duplicated local state is how the
  heart on the grid ends up disagreeing with the heart on the pack page.
- The toggle must feel instant: mutate the store first, persist after. On failure, roll the
  store back and show an error. Alpine's reactivity makes the revert a state change, not a
  DOM edit.
- The favorite button sits inside the tile, which is itself the copy target. **Stop
  propagation** or favoriting will also copy.
- Hydrate in **one** pass: `store.load()` reads every favorite id once into a `Set` on
  `alpine:init`. Tiles then bind `:class="$store.favorites.has(id) && 'is-fav'"`. Never an
  `isFavorite()` round-trip per tile.
- Alpine 3 reactivity does not track `Set` mutations reliably — reassign
  (`this.ids = new Set(this.ids)`) after add/remove, or hold the ids in a plain array or
  object. Mutating a `Set` in place and expecting the UI to update is a silent no-op.
- Never send favorite state to the server. Not even analytics.

## Acceptance

- [ ] Every tile shows a favorite toggle reflecting stored state on load
- [ ] Toggling on stores metadata and caches the blob
- [ ] Toggling off removes both
- [ ] Clicking the toggle does not trigger a copy
- [ ] State survives a page reload
- [ ] A failed write reverts the icon and shows an error
- [ ] Hydration issues one store read, not N
- [ ] Toggling on the grid updates the same sticker's heart on the pack page — one store,
      not two copies of the state

## Verify

```
make run
```

Favorite three stickers → reload → still marked. DevTools → Application → IndexedDB shows
three records in `favorites` and three in `blobs`. Network tab shows no request carrying
favorite data.

## Files

- `public/static/js/favorites.js` (data ops; framework-free)
- `public/static/js/app.js` (the `favorites` store)
- `pkg/ui/components/sticker.go`

## Hints

- Import direction is one-way: `app.js` → `favorites.js` → `db.js`. Nothing imports
  `app.js`.
- The reassign-the-Set gotcha above will cost you 20 minutes if you meet it cold. Use a
  plain object keyed by id if you would rather not think about it.
- Reuse the T11 container handler; branch on `closest('[data-fav-toggle]')` before the
  copy branch, and `$event.stopPropagation()` so favoriting does not also copy.
- A heart that fills instantly and reverts on failure reads far better than a spinner.
