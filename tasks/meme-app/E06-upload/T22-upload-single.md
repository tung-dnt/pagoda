# T22 · Upload POST: single image → Singles

| | |
|---|---|
| **Epic** | [E06 · Upload](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T21, T06 |
| **Unlocks** | T23 |
| **Spec** | §6.5 |
| **Status** | ☐ not started |

## Problem

Handle `POST /upload` for a **single image**: validate it, store it in the pending bucket,
and append a pending sticker row to the global "Singles" pack.

Single image first, zip second — this path is simpler and proves the whole pipeline
end-to-end before async processing enters the picture.

## Given

- Gate from T21, `Inspect` from T06, `Storage` from T05
- `pkg/handlers/files.go` — the reference multipart handler (**delete it once this works**,
  per T03)

## Constraints

- Request body capped at **100 MB** — enforced by middleware, before the handler allocates
  anything.
- Validate with `Inspect` **before** anything is written to R2. Never store then check.
- The image goes to `pending/{sticker_id}.{ext}` in the **private** bucket. It must not be
  publicly reachable before approval.
- The sticker is appended to the pack where `is_global = TRUE`. **No new pack is created.**
- Row is written with `status = 'pending'`, and `animated` from `Inspect` — never guessed
  from the extension.
- Order matters: write the object first, then the row. A row pointing at a missing object
  is a broken sticker; an orphan object is merely wasted bytes. If the row insert fails,
  delete the object.
- Success shows a clear "submitted for review" message. The uploader must not expect it to
  appear immediately.

## Acceptance

- [ ] A valid PNG upload creates one pending sticker in the "Singles" pack
- [ ] No new pack is created
- [ ] The object lands in the pending bucket and is **not** reachable via the public URL
- [ ] An invalid image is rejected with a reason, and nothing is written anywhere
- [ ] A 150 MB body is rejected by middleware, not by the handler
- [ ] An animated GIF is stored with `animated = true`
- [ ] The pending sticker is absent from browse and search

## Verify

```
make test && make run
```

Upload a PNG → check `psql` for a pending row on the global pack → try the public URL and
expect 403/404 → confirm it is absent from `/` and `/search`.

## Files

- `pkg/handlers/upload.go`, `pkg/handlers/upload_test.go`
- `pkg/postgres/queries/stickers.sql`
- `pkg/handlers/router.go` (body-limit middleware)

## Hints

- `echomw.BodyLimit("100M")` on the upload group only — not globally.
- Insert the row and get its id first to name the object key, then upload, then update the
  key. Or generate a UUID key up front and avoid the two-step entirely. Prefer the latter.
- The "not publicly reachable" check is the one that catches a misconfigured bucket. Do it
  by hand at least once.
