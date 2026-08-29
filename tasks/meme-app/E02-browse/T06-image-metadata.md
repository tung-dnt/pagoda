# T06 · Image sniffing + animation detection

| | |
|---|---|
| **Epic** | [E02 · Browse](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | None |
| **Unlocks** | T07, T22 |
| **Spec** | §6.3, §6.5 |
| **Status** | ☐ not started |

## Problem

Given an image's bytes, decide three things: is it an accepted format, what are its
dimensions, and **is it animated**. The `animated` flag drives which of the two copy paths
the client takes, so it must be determined server-side and stored — never guessed by JS.

This is the most self-contained, most testable task in the whole plan. Pure function,
bytes in, struct out.

## Given

- Go stdlib: `image`, `image/png`, `image/gif`, `image/jpeg`, `net/http.DetectContentType`
- WebP is not in the stdlib — parse its header manually, it is simple

## Constraints

- **Sniff content, never trust the extension.** A file named `cat.png` containing PHP must
  be rejected.
- Animation detection, by format:
  - **GIF** — more than one frame. `gif.DecodeAll` then `len(g.Image) > 1`
  - **PNG** — presence of an `acTL` chunk before the first `IDAT` (APNG)
  - **WebP** — RIFF container with an `ANIM` chunk
  - **JPEG** — never animated
- Accepted types are exactly: PNG, JPEG, GIF, WebP. Everything else is rejected with a
  reason string the uploader will see.
- Enforce dimensions: min 32×32, max 2048×2048.
- Must not read the whole file into memory twice — take an `io.ReadSeeker`.

## Acceptance

- [ ] `Inspect(r io.ReadSeeker) (Meta, error)` returns mime, width, height, animated
- [ ] A single-frame GIF reports `animated: false`; a multi-frame GIF reports `true`
- [ ] An APNG reports `animated: true`; a plain PNG reports `false`
- [ ] An animated WebP reports `true`; a static WebP reports `false`
- [ ] A `.png` whose bytes are a shell script is rejected
- [ ] A 16×16 image and a 4096×4096 image are both rejected with distinct reasons

## Verify

```
go test ./pkg/services/ -run TestInspect -v
```

Commit small fixture files under `pkg/services/testdata/`.

## Files

- `pkg/services/imagemeta.go`
- `pkg/services/imagemeta_test.go`
- `pkg/services/testdata/*` (8–10 tiny fixtures)

## Hints

- APNG: after the 8-byte PNG signature, walk chunks reading a 4-byte big-endian length
  and a 4-byte type. Stop at `IDAT`. If you saw `acTL`, it is animated.
- WebP: bytes 0–3 are `RIFF`, 8–11 are `WEBP`. Then walk chunks the same way looking for
  `ANIM`. `VP8X` with the animation bit set is the other tell.
- Build fixtures with ImageMagick or grab tiny ones from Wikipedia; keep every file under
  a few KB so the repo stays small.
