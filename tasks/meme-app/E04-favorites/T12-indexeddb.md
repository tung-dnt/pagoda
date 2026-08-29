# T12 · db.js: IndexedDB wrapper

| | |
|---|---|
| **Epic** | [E04 · Favorites](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | None |
| **Unlocks** | T13 |
| **Spec** | §6.4 |
| **Status** | ☐ not started |

## Problem

Write the storage layer for favorites: a small promise-based wrapper over IndexedDB with
three object stores. Everything about favorites is client-side — the server never learns
what anyone favorited.

| Store | Key | Holds |
|---|---|---|
| `favorites` | sticker id | sticker metadata + `favoritedAt` |
| `favorite_packs` | pack id | pack metadata + its sticker ids |
| `blobs` | object key | the image `Blob` + `bytes` + `lastUsed` |

## Given

- Nothing. This is a greenfield module with no dependencies on the rest of the app.
- Runs against `public/static/js/` — no bundler, plain ES module.
- **Framework-free**, like `clipboard.js`. Alpine is a UI layer, not a data layer; this
  module must not reference it. `app.js` imports this, never the reverse.

## Constraints

- Raw IndexedDB is event-based and unpleasant. Wrap every request in a Promise once, at
  the bottom, and never touch `onsuccess` above that layer.
- Version the schema and implement `onupgradeneeded` properly from v1. Getting this wrong
  means users with existing data hit an error you cannot reproduce.
- **Every read and write must survive being denied.** Private windows, Safari with site
  data blocked, and quota exhaustion all throw. A thrown IndexedDB error must degrade to
  "favorites unavailable", never a broken page.
- Blob storage is capped at **100 MB**. Enforce it here, in the storage layer, not in
  calling code — a cap that lives in the caller gets forgotten by the next caller.
- Eviction is least-recently-used by `lastUsed`, and evicts **blobs only**. Metadata is
  never evicted, so a favorite survives eviction and simply re-fetches on next view.

## Acceptance

- [ ] `open()`, `getFavorites()`, `addFavorite()`, `removeFavorite()`, `isFavorite()`
- [ ] `putBlob(key, blob)`, `getBlob(key)`, and `evictIfOver(limit)`
- [ ] Same API for the pack stores
- [ ] Exceeding 100 MB evicts oldest-used blobs until under the cap
- [ ] Evicting a blob leaves its favorite metadata intact
- [ ] Every method rejects cleanly (never throws synchronously) when IndexedDB is blocked

## Verify

Manual, in DevTools console:

```js
const db = await MemeDB.open()
await db.addFavorite({id: 1, name: 'test', url: '...'})
await db.getFavorites()
```

Then Application → IndexedDB and confirm the records. Repeat the whole flow in a private
window and confirm it degrades instead of exploding.

## Files

- `public/static/js/db.js`

## Hints

- One `promisify(request)` helper turns the whole API pleasant. Write it first.
- Track total blob bytes in a small meta record rather than summing every blob on each
  write — summing gets slow exactly when the cache is full.
- Test the cap by temporarily lowering it to 1 MB and favoriting a few stickers.
