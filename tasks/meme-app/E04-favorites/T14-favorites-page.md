# T14 · /favorites page + hydration

| | |
|---|---|
| **Epic** | [E04 · Favorites](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T13 |
| **Unlocks** | T15, T16 |
| **Spec** | §6.4 |
| **Status** | ☐ not started |

## Problem

`GET /favorites` returns an **empty shell**. The page is populated entirely from IndexedDB
by JS. This is the payoff task: it must work with the network switched off.

## Given

- `db.js`, `favorites.js`, and the `favorites` Alpine store from T13
- `pkg/ui/layouts/primary.go`

## Constraints

- The server returns no favorite data. It cannot — it does not have any.
- The page is one Alpine component (`x-data="favoritesPage"`) that loads from the store on
  `init()` and renders with `x-for`. Loading, empty, and populated are three states of one
  component, expressed with `x-show` — not three server-rendered branches.
- Render images from cached blobs via `URL.createObjectURL(blob)`. **Revoke every object
  URL** when its tile leaves the DOM, or the tab leaks memory for as long as it is open.
- If a blob was evicted but the metadata survives, re-fetch it and re-cache. Offline plus
  evicted is the one case that legitimately shows a placeholder.
- Copy works here exactly as on the pack page — reuse `clipboard.js` and the same
  container-delegation handler. Do not fork either.
- Empty state must be real UI, not a blank page.
- **State the limitation plainly on this page**: favorites live in this browser, clearing
  site data loses them, cross-device sync is phase 2. Spec §6.4 requires the UI to say so
  rather than let a user discover it the hard way.

## Acceptance

- [ ] `/favorites` renders every favorited sticker from IndexedDB
- [ ] With the network cut **mid-session**, the page still renders every sticker with
      zero image requests in the Network tab
- [ ] Copy works offline for cached static stickers
- [ ] Object URLs are revoked — heap does not grow across repeated navigations
- [ ] An evicted-blob favorite re-fetches when online, placeholders when offline
- [ ] Empty state renders when nothing is favorited
- [ ] The per-browser limitation is stated on the page

## Verify

1. Favorite several stickers
2. Open `/favorites` while online
3. DevTools → Network → **Offline**
4. Navigate away and back (client-side) → stickers still render, no image requests
5. Copy a static one → still works
6. Performance → two heap snapshots across navigations, confirm no growth

**A cold reload while offline is out of scope for this task** — the HTML shell still comes
from the server. [T16](./T16-service-worker.md) adds that and re-runs this verification.

## Files

- `pkg/ui/pages/favorites.go`
- `pkg/handlers/stickers.go`
- `public/static/js/favorites.js`
- `pkg/routenames/names.go`

## Hints

- The offline test is the acceptance criterion that actually proves "local-first". Do it
  before calling this done.
- Keep an id→objectURL map and revoke on replace and on `pagehide`.
