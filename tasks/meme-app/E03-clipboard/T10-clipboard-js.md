# T10 · clipboard.js: the two copy paths

| | |
|---|---|
| **Epic** | [E03 · Clipboard](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T01, T09 |
| **Unlocks** | T11 |
| **Spec** | §6.3 |
| **Status** | ☐ not started |

## Problem

Write the client module that puts a sticker on the clipboard. Two paths, chosen by
`data-animated`:

| Sticker | Mechanism | Pastes as |
|---|---|---|
| Static | `ClipboardItem({'image/png': blob})` | an image |
| Animated | `clipboard.writeText(url)` | a link that unfurls |

Use the exact code shape that worked in T01. Do not re-derive it.

## Given

- `docs/spike-clipboard.md` from T01 — findings and the working snippet
- Tiles from T09 carrying `data-sticker-url` and `data-animated`
- `public/static/` — currently has no JS at all; this is the first module

## Constraints

- **No bundler.** Plain ES module, served from `public/static/js/`, loaded with
  `<script type="module">`. No npm, no build step for JS.
- **Framework-free.** This module must not reference Alpine. Alpine owns UI state; this
  owns platform work. `app.js` imports it, never the reverse — that keeps it testable in
  isolation and survives Alpine being dropped.
- **The Safari rule:** construct `ClipboardItem` with a *Promise*, synchronously in the
  click handler. Never `await` the fetch before calling `clipboard.write`.
- Convert non-PNG static images to PNG via `<canvas>` before writing — `image/png` is the
  only universally supported clipboard image type.
- `fetch` the image with `mode: 'cors'`; a tainted canvas throws on `toBlob`. If the public
  bucket's CORS is wrong this is where you find out.
- Absolute URL for the animated path — a relative URL pasted into a chat app is useless.
- Every outcome shows a toast: **"Copied image"** or **"Copied link"**, and a distinct
  error toast on failure. Silent failure is forbidden by spec §8.
- Graceful degradation: if `navigator.clipboard` is missing (insecure context, old
  browser), show a message with the URL rather than throwing.

## Acceptance

- [ ] Clicking a static sticker copies a PNG that pastes as an image into Slack/Discord
- [ ] Clicking an animated sticker copies its absolute URL
- [ ] The toast always names which of the two happened
- [ ] A failed copy shows an error toast, never nothing
- [ ] Works in Chrome, Firefox, Safari desktop, and iOS Safari
- [ ] Zero console errors on a successful copy

## Verify

Manual, and this one genuinely needs the phone:

1. Serve over HTTPS (tunnel is fine)
2. Copy a static sticker → paste into iOS Messages → an image appears
3. Copy an animated sticker → paste → the animation plays
4. Repeat in desktop Chrome, Firefox, Safari

## Files

- `public/static/js/clipboard.js`
- `pkg/ui/components/head.go` (script tag)
- `pkg/ui/components/alerts.go` (reuse for the toast if it fits)

## Hints

- `ui.StaticFile()` appends a cache-busting key — use it for the script URL, or you will
  spend an hour debugging a cached module.
- Keep this file free of favorites logic. T13 imports it; it must not import T13.
- If CORS bites, fix it on the bucket, not with a proxy — routing bytes through Go is on
  the Never list.
