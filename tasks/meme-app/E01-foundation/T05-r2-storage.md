# T05 · R2 storage service

| | |
|---|---|
| **Epic** | [E01 · Foundation](./README.md) |
| **Size** | M — ~3 h |
| **Depends on** | None |
| **Unlocks** | T07, T21, T25 |
| **Spec** | §6.7 |
| **Status** | ☐ not started |

## Problem

Build `services.Storage`: a thin wrapper over the S3 API pointed at Cloudflare R2, with
two buckets — one private for pending uploads, one public for approved stickers.

Two buckets, not one, because R2 public access is per-bucket. Pending content must not be
reachable by guessing a URL.

## Given

- `pkg/services/container.go` — add the field and wire it in `NewContainer`
- `config/config.go` + `config/config.yaml` — the config pattern
- `aws-sdk-go-v2` (new dependency; already approved for this task only)

## Constraints

- **Credentials come from env, never `config.yaml`.** The file is committed.
- R2 needs `region: "auto"` and a custom endpoint
  (`https://<account>.r2.cloudflarestorage.com`). It is S3-compatible, not S3.
- Approved objects get `Cache-Control: public, max-age=31536000, immutable`. Keys are
  immutable — a change means a new key, never an overwrite.
- `PublicURL(key)` returns the public bucket's custom-domain URL. **Image bytes must never
  pass through the Go app** — that is the reason R2 was chosen. No proxy handler.
- Presigned GET URLs for admin preview: 15 min TTL, private bucket only.
- Local dev must work without real R2 credentials — see Hints.

## Acceptance

- [ ] `Put(ctx, key, r, contentType)` writes to the pending bucket
- [ ] `Promote(ctx, key)` copies pending → public and deletes the pending object
- [ ] `Delete(ctx, key)` removes from pending
- [ ] `PublicURL(key)` returns a URL on the public base, no signing
- [ ] `PresignedGet(ctx, key)` returns a URL that expires in 15 min
- [ ] Missing credentials produce a clear startup error, not a nil-pointer panic later
- [ ] `c.Storage` is available to handlers via the container

## Verify

```
make build
make test          # unit tests against a fake/in-memory implementation
```

Manual: put an object, promote it, fetch the public URL with `curl -I` and confirm the
`Cache-Control` header.

## Files

- `pkg/services/storage.go`, `pkg/services/storage_test.go`
- `pkg/services/container.go`
- `config/config.go`, `config/config.yaml`

## Hints

- Define a small interface (`Put`/`Promote`/`Delete`/`PublicURL`/`PresignedGet`) so tests
  use a memory fake and no test ever needs network.
- For local dev, either point at a real R2 dev bucket or back the fake with `afero` on
  `uploads/` — the existing `c.Files` already gives you that.
- `Promote` is `CopyObject` then `DeleteObject`. There is no atomic move; if the delete
  fails after the copy, log it and move on — a leaked pending object is harmless.
