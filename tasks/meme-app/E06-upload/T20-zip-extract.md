# T20 · Safe zip extraction service

| | |
|---|---|
| **Epic** | [E06 · Upload](./README.md) |
| **Size** | M — ~3 h |
| **Depends on** | T06 |
| **Unlocks** | T23 |
| **Spec** | §6.5 |
| **Status** | ☐ not started |

## Problem

Given an uploaded zip, produce a list of valid images — safely. This is the highest-risk
code in the app: it processes attacker-controlled archives on a public site.

Write it as a pure service with no HTTP and no storage dependencies, so it can be tested
adversarially in isolation.

## Given

- Go stdlib `archive/zip`
- `Inspect` from T06

## Constraints — every one of these is a real attack

| Guard | Limit | Attack |
|---|---|---|
| Entry count | 200 | Millions of tiny files |
| Total uncompressed | 250 MB | Decompression bomb |
| Per-image | 5 MB | Memory exhaustion |
| Compression ratio | reject > 100:1 | Nested bomb |
| Path | `filepath.Base` only | `../../etc/passwd` traversal |
| Type | content sniffing | `.png` containing a script |

- **Never trust `zip.File.UncompressedSize64`** — it is attacker-controlled metadata. Cap
  the actual read with `io.LimitReader` and compare against what you really read.
- Discard entry paths entirely. Only `filepath.Base(name)` is used. Directory entries,
  symlinks, and dotfiles are skipped.
- Pack name derivation: the single top-level folder name; if entries sit at the zip root,
  the zip's filename. An explicit form name overrides both.
- A failing entry is **skipped with a reason**, not fatal. The submission succeeds if at
  least one image is valid — a contributor should not lose 19 good stickers to 1 bad one.
- Stream entries; do not read the whole archive into memory.

## Acceptance

- [ ] `Extract(r io.ReaderAt, size int64) (Result, error)` returns valid images + per-entry
      skip reasons + the derived pack name
- [ ] A zip bomb (>100:1) is rejected without exhausting memory
- [ ] `../../etc/passwd` entries write nothing outside the target
- [ ] A `.png` containing a shell script is skipped with a reason
- [ ] A 201-entry zip is rejected
- [ ] A zip with 19 good and 1 bad image yields 19 images and 1 reason
- [ ] A zip with zero valid images returns a clear error

## Verify

```
go test ./pkg/services/ -run TestExtract -v
```

Every constraint above needs its own test with a committed fixture. This is the one place
in the plan where the tests matter more than the implementation.

## Files

- `pkg/services/zip.go`, `pkg/services/zip_test.go`
- `pkg/services/testdata/zips/*`

## Hints

- Build the bomb fixture with a script committed alongside it, not by hand.
- Track running totals across entries — the 250 MB cap is cumulative, and per-entry checks
  alone will not catch 200 files of 2 MB each.
- Wrap each `rc` in `io.LimitReader(rc, perFileMax+1)` and treat reading `perFileMax+1`
  bytes as "too big".
