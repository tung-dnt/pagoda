# T01 · Clipboard spike on iOS Safari

| | |
|---|---|
| **Epic** | [E01 · Foundation](./README.md) |
| **Size** | S — ~1 h |
| **Depends on** | None |
| **Unlocks** | T10 |
| **Spec** | §6.3 |
| **Status** | ☐ not started |

## Problem

Before building any UI, prove the core interaction actually works on the platform most
likely to break it. Write one throwaway static HTML page that copies a static PNG and an
animated GIF to the clipboard, and test it on a **real iPhone** over HTTPS.

This is task 1 of 28 on purpose. If iOS Safari can't do this, the product changes shape,
and you want to learn that today rather than in week five.

## Given

- Nothing in the repo yet. This is a standalone file, deliberately outside the app.
- Two sample images you can host anywhere reachable over HTTPS (one PNG, one GIF).

## Constraints

- **Secure context required.** `navigator.clipboard` is undefined on plain `http://` on a
  phone. Use a tunnel (`cloudflared tunnel --url http://localhost:8000`) or a static host.
  `localhost` is exempt on desktop but that is exactly the case that will mislead you.
- **The Safari trap.** `ClipboardItem` must be constructed with a *Promise* for the blob,
  synchronously inside the click handler. This works:

  ```js
  navigator.clipboard.write([new ClipboardItem({ 'image/png': fetchAsPng(url) })])
  ```

  This silently fails on Safari, because the `await` loses the user-gesture context:

  ```js
  const blob = await fetchAsPng(url)                       // ✗ gesture is gone
  navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })])
  ```
- The remote image must be served with permissive CORS, or `canvas.toBlob()` throws a
  security error on a tainted canvas. An `<img>` that renders fine can still taint.
- Only `image/png` is reliably accepted. Do not spend time trying to make GIF work as an
  image — confirming it fails is part of the deliverable.

## Acceptance

- [ ] Tapping "copy static" on a real iPhone puts an image on the clipboard that pastes
      into iOS Messages as a picture
- [ ] Tapping "copy animated" puts a URL on the clipboard that pastes as text and unfurls
- [ ] The same page behaves correctly in desktop Chrome, Firefox, and Safari
- [ ] Findings written to `docs/spike-clipboard.md` — what worked, what silently failed,
      and any per-browser quirk worth knowing in T10

## Verify

Manual only. Open on a physical iPhone over HTTPS, paste into Messages and into WhatsApp
or Zalo. A simulator is not sufficient — clipboard behaviour differs.

## Files

- `docs/spike-clipboard.md`
- a scratch HTML file (do **not** commit it to `public/`)

## Hints

- Convert non-PNG to PNG by drawing to a `<canvas>` and calling `toBlob(cb, 'image/png')`.
- Detect the failure mode by wrapping the write in `try/catch` **and** logging the
  resolved promise — Safari sometimes rejects without throwing synchronously.
- If it works, note the exact code shape that worked. T10 copies it verbatim.
