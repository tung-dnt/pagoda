# Implementation Plan: Meme Sticker App — Phase 1

Spec: [`docs/spec/phase-1.md`](../docs/spec/phase-1.md) · Intent: [`docs/intent/meme-app.md`](../docs/intent/meme-app.md)
Board: [`tasks/meme-app/`](./meme-app/) · Checklist: [`tasks/todo.md`](./todo.md)

## Overview

A public sticker site on the existing Pagoda stack (Go 1.27, Echo, gomponents, HTMX,
pgx/sqlc). Visitors browse packs, copy a sticker in one click, and favorite stickers with
no account. Contributors upload behind an invite code; a single admin approves. 29 tasks
across 8 epics, roughly 62 hours.

## Architecture decisions

- **Approval status lives on `stickers`, not `packs`.** A pack is a stateless container and
  becomes visible when it has ≥1 approved sticker. Lets the admin reject 2 images out of a
  20-image zip, and makes the global "Singles" pack an ordinary case rather than a special
  one. Visibility is derived from a join, never stored, so the two cannot drift.
- **Two R2 buckets.** R2 public access is per-bucket, so pending content needs its own
  private bucket. Approval copies the object across and deletes the original.
- **Image bytes never pass through the Go app.** Approved images are served from the public
  bucket's custom domain. This is the reason R2 was chosen; a proxy handler would negate it.
- **`animated` is decided server-side at upload** and stored on the row. The client never
  guesses which copy path to take from a file extension.
- **No identity mechanism at all in phase 1.** Favorites are pure IndexedDB — no device id,
  no cookie, no server-side favorites table. Passkeys are phase 2. Device fingerprinting was
  dropped entirely: it cannot work across a user's devices and collides on a public site.
- **Offline is a service worker, not just IndexedDB.** Blobs in IndexedDB make images
  local; a cold reload still needs the HTML shell, CSS and JS. T16 precaches the shell.
  Sticker images are deliberately *not* in the SW cache — they are already in IndexedDB
  and double-storing them would defeat T12's eviction cap.
- **Alpine is pinned to an exact version (3.16.3), not `3.x.x`.** A floating range behind
  a stable service-worker cache key serves changing content from the same URL. The CDN
  dependency itself is retained by choice; vendoring locally remains the safer option.
- **No JS bundler.** Plain ES modules from `public/static/js/`, plus **Alpine.js 3** and
  HTMX loaded from unpkg in `pkg/ui/components/head.go`. Tailwind still builds CSS.
- **Alpine owns UI state; plain modules own data and platform work.** `db.js` (IndexedDB)
  and `clipboard.js` (fetch, canvas, clipboard) stay framework-free and testable in
  isolation; `app.js` holds every `Alpine.data()` / `Alpine.store()` registration. Import
  direction is one-way — nothing imports `app.js`.
- **One `Alpine.store('favorites')`, not per-component state.** The grid, pack page, and
  favorites page read the same store, so a toggle anywhere updates everywhere.

## Dependency graph

```
T01 clipboard spike ─────────────────────────────────┐
T02 strip auth ──→ T03 strip demos ──┐               │
T04 schema ──┬───────────────────────┤               │
             │                       ↓               ↓
T05 storage ─┼──→ T07 seed ──→ T08 grid ──→ T09 detail ──→ T10 clipboard.js
             │         ↑                        │              ↓
T06 imagemeta┴─────────┘                        │         T11 copy wiring
                                                │              ↓
                                    T17 tags ───┤         T13 fav toggle ←── T12 db.js
                                        ↓       │              ↓
                                    T18 search  │         T14 fav page
                                        ↓       │              ↓
                                    T19 filter  └────────→ T15 fav pack
                                                               ↓
                                                          T16 service worker

T06 ──→ T20 zip extract ──┐
T05 ──→ T21 invite gate ──┴──→ T22 upload single ──→ T23 upload zip
                                                          ↓
                                    T24 queue ──→ T25 approve/reject ──→ T26 bulk
                                                                            ↓
T27 CI (independent) ──────────────────────→ T28 browser checklist ──→ T29 deploy
```

## Phases and checkpoints

| Epic | Tasks | Budget | Checkpoint |
|---|---|---|---|
| [E01 Foundation](./meme-app/E01-foundation/) | T01–T05 | ~9 h | Migrations reversible, storage round-trips, clipboard proven on iOS |
| [E02 Browse](./meme-app/E02-browse/) | T06–T09 | ~8 h | Real stickers visible; pending content absent from HTML source |
| [E03 Clipboard](./meme-app/E03-clipboard/) | T10–T11 | ~4 h | **Core loop works** — copy and paste into a chat app from a phone |
| [E04 Favorites](./meme-app/E04-favorites/) | T12–T16 | ~12 h | Cold reload while offline still renders favorites |
| [E05 Discovery](./meme-app/E05-discovery/) | T17–T19 | ~6 h | Search and tag filter return approved content only |
| [E06 Upload](./meme-app/E06-upload/) | T20–T23 | ~10 h | Zip of 20 lands as 20 pending stickers; every abuse test passes |
| [E07 Moderation](./meme-app/E07-moderation/) | T24–T26 | ~7 h | Approving the first sticker makes a pack appear publicly |
| [E08 Ship](./meme-app/E08-ship/) | T27–T29 | ~6 h | All 16 spec success criteria verified against production |

Each epic README carries its own checkpoint. **Every task ends with the repo green**
(`make test && make build` exits 0). A task that cannot was too big — split it.

## Sequencing rationale

- **T01 is first on purpose.** The clipboard spike is the highest-risk unknown in the
  project and the cheapest to test. If iOS Safari cannot do it, the product changes shape,
  and that is worth knowing on day one rather than in week five.
- **E01 is deliberately horizontal**, which the vertical-slicing rule normally forbids.
  Schema and storage genuinely underpin every slice, and the honest alternative — building
  them piecemeal inside four different feature tasks — is worse. E02 onward are proper
  vertical slices.
- **E03 is placed before favorites and search** so the core interaction is proven working
  end-to-end at hour ~21, not at the end.
- **T20 (zip extraction) is a pure service with no HTTP or storage dependencies**, so the
  adversarial tests that matter most can be written in isolation.

## Risks and mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| iOS Safari clipboard doesn't work as specced | **High** — core interaction | T01 spikes it standalone before any UI exists |
| Client JS has no automated tests | **High** — silent regressions | T28's written checklist across 4 browsers; deliberate, time-driven cut |
| R2 CORS misconfigured — `fetch` fails where `<img>` works | **High** — copy and blob caching both break, production-only | Explicit CORS verification in T29, tested with `fetch` not `<img>` |
| Zip handling is attacker-facing on a public site | **High** — RCE / DoS surface | T20 is pure and isolated; 6 named guards, each with its own test |
| Approve promotes object but DB update fails | Medium — orphan object | Promote first, then flip the row. Orphan bytes are harmless; a visible broken image is not |
| Global "Singles" pack grows unbounded | Medium — unbounded blob download | Pack-favorite control not rendered when `is_global`; on the spec's Never list |
| IndexedDB denied (private window, Safari ITP) | Medium — favorites vanish | T12 degrades cleanly; T14 states the limitation in the UI |
| Search unindexed as library grows | Low — deferred by choice | `ILIKE` now; `pg_trgm` when a query is measurably slow (Open Q2) |
| unpkg unreachable on a first visit | Medium — SW install fails, no offline | Handle a failed install gracefully; vendoring would remove this entirely |
| Service worker serves stale assets after deploy | Medium — hard to debug | Network-first for HTML, versioned cache, no unconditional `skipWaiting()` |

## Parallelization

Mostly a solo, sequential build, but if a second session is available:

- **Independent, start any time:** T01, T02, T04, T05, T06, T12, T27
- **Must be sequential:** T04 → everything touching the schema; T22 → T23; T24 → T25 → T26
- **Needs the contract first:** T10/T13/T15 all read the tile data attributes defined in
  T09. Fix those attribute names in T09 and the three client tasks can then proceed apart.

## Open questions

Carried from the spec, none blocking:

1. **Go version** — CI pins 1.25, `go.mod` says 1.27.0. Resolved by T27.
2. **`pg_trgm`** — add the GIN index now, or wait for a measured slow query?
3. **Invite code rotation** — config-only means rotation needs a redeploy. Acceptable?
4. **Public bucket domain** — custom domain, or `r2.dev` for phase 1?
5. **Rejected stickers** — keep the row for audit, or purge after N days?
