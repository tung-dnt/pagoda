# Intent — Meme Sticker App (Phase 1)

Confirmed via `interview-me` on 2026-08-25. This is a statement of *intent*, not a spec.
Downstream: `spec-driven-development`.

## Core

| | |
|---|---|
| **Outcome** | A public sticker site — fast browsing, one-click copy, favorites that persist with no signup screen. |
| **User** | Public-facing, but built for the author and their circle; they are the ones who must actually reach for it. |
| **Why now** | The Pagoda boilerplate is standing and deserves a real product on it. |
| **Success** | Live, and the author + friends use it instead of scrolling a camera roll. |
| **Constraint** | Author's own time. Evenings/weekends. Every subsystem must earn its place. |

## Stack (already fixed by the repo)

Go 1.27 · Echo · **maragu.dev/gomponents** (not templ) · HTMX · pgx + sqlc · `cmd/web` + `cmd/admin`.

Views are pure Go functions returning `Node`, composed via `pkg/ui/components` and `pkg/ui/layouts`,
with a render cache (`cache.SetIfNotExists`).

**No build step, no bundler.** The client-side favorites layer (IndexedDB, clipboard writes, passkey
call) is hand-written vanilla JS served as a static asset. It is the one part of phase 1 that does
not get Go's type checking — budget review effort accordingly.

## Phase 1 scope

- Browse / surf stickers; **packs** are the grouping unit
- Copy to clipboard:
  - static → clipboard as PNG
  - animated → clipboard as **link**
- Favorite a single sticker **or** an entire pack
- Favorites are **local-first** in IndexedDB; anonymous device ID by default
- **Passkey** as opt-in for cross-device favorites
- Upload page gated by an **invite code** — single image or zip-of-folder
- Admin approval queue; single admin account

## Out of scope (phase 1)

Open public uploads · NSFW detection · takedown flow · spam defense · user accounts/profiles ·
comments/reactions · sticker editing or generation · multiple admins · analytics · monetization ·
native/mobile apps.

## Decisions that overrode the original ask

1. **Device fingerprinting is dropped.** It cannot do what was asked. Fingerprints are unrelated
   across a user's Apple devices by construction (different GPU, canvas, screen metrics, audio
   stack), and Safari actively randomizes them. Apple exposes no cross-device identifier to web
   content, deliberately. Worse, on a *public* app fingerprints **collide** — two stock iPhones on
   the same iOS version can hash identically, serving stranger A's favorites to stranger B. That is
   a correctness bug, not just a privacy one.
   → Replacement: anonymous device-local ID by default, **passkey (WebAuthn)** offered only when a
   user wants favorites on a second device. Syncs via iCloud Keychain, one Face ID tap, no email.

2. **Copy-to-clipboard splits in two.** Browsers can only write PNG to the clipboard via
   `ClipboardItem`. An animated GIF/WebP copied as an image pastes as a dead first frame. This is a
   platform limit with no workaround. Mixed libraries therefore need two copy paths and a UI cue for
   which one the user is getting.

3. **Uploads cut to invite-only.** Open-internet uploads were the original ask, but they are the
   most expensive option available and buy nothing on day one when nobody knows the site exists.
   Given the time constraint, phase 1 gates uploads behind a shared invite code; opening to the
   public — with rate limiting, zip-bomb guards, NSFW review, and a takedown path — is phase 2.

## Resolved in the spec

See `docs/spec/phase-1.md`.

- One zip = one pack; the top-level folder name becomes the pack name (form name overrides).
- A **single image** goes into one shared global "Singles" pack, not a pack of its own.
- Consequently, approval status lives on **stickers**, not packs: a pack is visible once it
  has ≥1 approved sticker, and the admin can reject part of a zip. The global pack cannot be
  favorited as a whole (it grows without bound).
- Tags **and** search both ship in phase 1.
- Storage is **Cloudflare R2**, two buckets (pending private, approved public).
- **Passkeys are deferred to phase 2.** Consequence: phase 1 has no identity mechanism at all —
  no device ID, no cookie, no server-side favorites table. Favorites are purely IndexedDB.
