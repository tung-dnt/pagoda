# T11 · Wire copy into the sticker tile

| | |
|---|---|
| **Epic** | [E03 · Clipboard](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T10 |
| **Unlocks** | T13 |
| **Spec** | §6.3 |
| **Status** | ☐ not started |

## Problem

Connect `clipboard.js` to the tiles rendered by T09, so a tap anywhere on a sticker copies
it. Make the affordance obvious and give immediate feedback.

## Given

- `clipboard.js` from T10
- `StickerTile` from T09
- **Alpine.js 3**, already loaded in `pkg/ui/components/head.go`

## Constraints

- **One delegated handler on the grid container**, not one per tile. In Alpine that is
  `x-data="stickerGrid"` plus `@click="onClick($event)"` on the container, resolving the
  tile with `e.target.closest('[data-sticker-id]')`. Putting `@click` on each tile means
  Alpine registers hundreds of listeners on a large pack.
- Register `stickerGrid` via `Alpine.data()` in `app.js` inside an `alpine:init` listener.
  Alpine is loaded `defer`, so **our module tag must come before the Alpine script tag** in
  `head.go` or the component is registered too late and `x-data` silently resolves to
  nothing.
- Never inline a JS object literal in `x-data`. The view says `x-data="stickerGrid"`; the
  behaviour lives in `app.js`.
- The tile must be a real `<button>` (or have `role="button"` and be keyboard-focusable).
  A clickable `<div>` is not reachable by keyboard.
- Visible state change on copy — a brief checkmark or ring on the tile itself, in addition
  to the toast. The toast alone is easy to miss on mobile. This is Alpine's job: a
  `copiedId` in component state and `:class` on the tile, not manual `classList` calls.
- An animated sticker should hint that it copies a link, before the click. A small badge is
  enough. Users should not have to discover the difference by pasting.
- Copy must still work for tiles added to the DOM later. Alpine 3 auto-initializes new
  nodes, and container-level delegation covers it regardless.

## Acceptance

- [ ] Clicking any tile copies via the correct path
- [ ] Exactly one click listener is attached per grid, verified in DevTools
- [ ] `x-data` resolves — no "Alpine Expression Error" in the console
- [ ] Tab-to-focus then Enter copies
- [ ] Animated tiles carry a visible badge
- [ ] The tile itself shows feedback, not just the toast

## Verify

```
make run
```

DevTools → Elements → Event Listeners on the grid: exactly one `click`. Tab through the
grid with the keyboard. Test on a phone, where hover states do not exist.

## Files

- `public/static/js/app.js` (new — Alpine registrations)
- `pkg/ui/components/sticker.go`
- `pkg/ui/components/head.go` (script tag ordering)

## Hints

- `e.target.closest('[data-sticker-id]')` is the whole delegation trick.
- gomponents needs no Alpine helper — directives are ordinary attributes:
  `Attr("x-data", "stickerGrid")`, `Attr("@click", "onClick($event)")`.
- If `x-data` appears to do nothing, check script order in `head.go` first. It is almost
  always that.
- Reuse the existing DaisyUI toast/alert classes from `pkg/ui/components/alerts.go` rather
  than inventing a new one.
