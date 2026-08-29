# T16 · Service worker: offline shell

| | |
|---|---|
| **Epic** | [E04 · Favorites](./README.md) |
| **Size** | M — ~3 h |
| **Depends on** | T14, T15 |
| **Unlocks** | T28, T29 |
| **Spec** | §6.4 |
| **Status** | ☐ not started |

## Problem

T14 made favorites render from IndexedDB, but a **cold reload while offline** still fails —
the HTML shell, the CSS, and the JS all come from the network. Add a service worker that
precaches the shell so `/favorites` works from a dead-cold start with no connection.

This is what makes the app usable in the situation it was built for: opening it mid-chat on
a phone with bad signal.

## Given

- `/favorites` and the pack pages from T14 / T15
- `public/static/` served by `router.go` via `echo.MustSubFS(files.Static, "static")`
- Alpine pinned at `3.16.3` and htmx at `2.0.0`, both from unpkg

## Constraints

- **Scope is the trap.** A worker served from `/static/sw.js` gets scope `/static/` and
  cannot control `/favorites`. Either serve it from the root path `/sw.js` with its own
  Echo route, or send a `Service-Worker-Allowed: /` header. Do the former; it is simpler
  and harder to get wrong.
- **Never cache sticker images.** They already live in IndexedDB as blobs. Caching them
  again doubles storage against the same quota and makes the T12 eviction cap meaningless.
  Precache the shell only: HTML, CSS, JS, fonts, favicon.
- The two unpkg scripts must be precached, or Alpine and htmx are missing offline and the
  page renders nothing. unpkg sends `Access-Control-Allow-Origin: *`, so a normal `cors`
  fetch works and gives a real (non-opaque) response. **Install will fail if unpkg is
  unreachable on that first visit** — this is the cost of the CDN over vendoring; handle a
  failed install without breaking the page.
- Versioned cache name (`shell-v1`), and delete every other cache on `activate`. Without
  this, users accumulate stale caches forever.
- HTML: **network-first, falling back to cache.** Cache-first on HTML ships stale pages
  after every deploy and is very hard to debug.
- Static assets: cache-first. They are already cache-busted by `ui.StaticFile()`.
- **Do not call `skipWaiting()` unconditionally.** Activating a new worker under a page
  running the old assets gives version-skew bugs that reproduce for nobody. Let it activate
  on next navigation, or gate it behind an explicit user action.
- **Register only when `app.environment != "local"`**, or add a documented dev escape
  hatch. A service worker serving stale JS during development will cost you an evening
  before you realise what is happening.

## Acceptance

- [ ] `GET /sw.js` is served from the root path with the correct content type
- [ ] The worker registers and its scope is `/`
- [ ] With DevTools **Offline**, a hard reload of `/favorites` renders every favorited
      sticker
- [ ] With DevTools **Offline**, a hard reload of `/` renders the shell (an empty-state or
      cached grid, not the browser's error page)
- [ ] No sticker image URL appears in the Cache Storage entries
- [ ] Bumping the cache version removes the old cache on activate
- [ ] A deploy of new CSS/JS reaches users without a manual cache clear
- [ ] It does not register in local development, or the escape hatch is documented

## Verify

```
make build && make run
```

1. Load the site, favorite several stickers
2. DevTools → Application → Service Workers: confirm activated, scope `/`
3. Application → Cache Storage: shell assets present, **no** sticker images
4. Network → **Offline**, then **hard reload** `/favorites` → stickers render
5. Bump the cache version, reload twice, confirm the old cache is gone
6. Repeat step 4 on a real iPhone — iOS evicts service workers more aggressively

## Files

- `public/sw.js` (root-scope, not under `static/`)
- `pkg/handlers/router.go` (the `/sw.js` route)
- `public/static/js/app.js` (registration)

## Hints

- Precache list: `/`, `/favorites`, `main.css`, the three JS modules, the two unpkg URLs,
  `favicon.png`.
- `caches.match(event.request)` ignores query strings only if you pass
  `{ignoreSearch: true}` — relevant because `ui.StaticFile()` appends a cache key.
- DevTools → Application → Service Workers → **Update on reload** makes iteration bearable.
- If offline reload shows the browser error page, the scope is wrong. Check it first.
