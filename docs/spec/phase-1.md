# Spec: Meme Sticker App — Phase 1

Status: **draft, awaiting review**
Consumes: [`docs/intent/meme-app.md`](../intent/meme-app.md)
Next: `planning-and-task-breakdown` → `tasks/plan.md` + `tasks/todo.md`

---

## 1. Objective

A public sticker site where browsing is fast, copying a sticker is one click, and favorites
persist with no signup screen.

**Users**
- **Visitor** (anonymous, no session) — browses packs, searches, copies, favorites.
- **Contributor** (holds the invite code) — uploads a single image or a zip.
- **Admin** (exactly one account) — approves or rejects submissions.

**Success:** the site is live and the author + friends reach for it instead of scrolling a
camera roll.

**Binding constraint:** author's own time. Every subsystem must earn its place.

### User stories

| # | As a… | I want… | So that… |
|---|---|---|---|
| U1 | Visitor | to browse a grid of sticker packs | I can find something funny |
| U2 | Visitor | to open a pack and see its stickers | I can pick one |
| U3 | Visitor | to click a sticker and have it on my clipboard | I can paste it into a chat |
| U4 | Visitor | to favorite one sticker or a whole pack | I can find it again fast |
| U5 | Visitor | my favorites to load instantly and work offline | the app feels immediate |
| U6 | Visitor | to search by name and filter by tag | I can find a specific sticker |
| U7 | Contributor | to upload a zip or single image with the invite code | I can add content without server access |
| U8 | Admin | to review pending submissions and approve/reject | nothing embarrassing goes public |

---

## 2. Tech Stack

Fixed by the existing repo (module `github.com/tung-dnt/meme-app`, a Pagoda-derived boilerplate).

| Concern | Choice |
|---|---|
| Language | Go 1.27.0 |
| HTTP | Echo v4.14 |
| Views | `maragu.dev/gomponents` v1.3 — **not** templ |
| Interactivity | HTMX + **Alpine.js 3** (both from unpkg CDN, see `pkg/ui/components/head.go`) |
| CSS | Tailwind CLI + DaisyUI (`make css`) |
| DB | PostgreSQL 17, `pgx/v5`, `sqlc` v1.31 (`emit_interface`, `emit_empty_slices`) |
| Migrations | `golang-migrate` v4, embedded, applied on startup |
| Local files | `afero` via `c.Files` (retained for scratch/temp only) |
| Object storage | **Cloudflare R2** via `aws-sdk-go-v2` S3 client (new dependency) |
| Tasks | `backlite` (SQLite) — used for async zip extraction |
| Binaries | `cmd/web` (site + admin UI), `cmd/admin` (create admin user) |

**No JS bundler.** Client-side code is plain ES modules served from `public/static/js/`,
plus Alpine directives in the gomponents markup. This is the only part of phase 1 without
compile-time type checking; it gets manual browser verification instead of automated tests
(see §7).

**Division of labour on the client.** Alpine owns *UI state* — toggles, badges, toasts,
progress, optimistic updates. Plain ES modules own *data and platform work* — IndexedDB,
`fetch`, blob handling, clipboard writes. Alpine is not a data layer; `db.js` and
`clipboard.js` must stay framework-free so they remain testable in isolation and reusable
if Alpine is ever dropped.

**Alpine loads deferred from a CDN**, which imposes two rules:
- Any `Alpine.data()` / `Alpine.store()` registration must happen inside an
  `alpine:init` listener, and our module tag must sit **before** the Alpine script tag in
  `head.go`, or Alpine boots before the components are registered.
- Alpine 3 auto-initializes nodes added to the DOM, so HTMX-swapped fragments containing
  `x-data` work without a manual re-init call.

---

## 3. Commands

```
make install                      # Go modules + Tailwind CLI + DaisyUI
make db-up                        # Start PostgreSQL via Docker
make db-down                      # Stop PostgreSQL

make migrate-new name=add_stickers # Create a migration pair
make migrate-up                   # Apply pending migrations
make migrate-down                 # Roll back the most recent migration
make migrate-force version=N      # Clear a dirty migration state

make sqlc-gen                     # Regenerate pkg/postgres/db from queries
make sqlc-vet                     # Lint the SQL queries

make admin email=admin@home.local # Create the single admin user

make css                          # Build + minify Tailwind CSS
make build                        # css + compile the binary
make run                          # Run the application
make watch                        # Run with air, rebuild on change
make test                         # All tests (requires make db-up)
```

**Definition of "green":** `make sqlc-vet && make test && make build` all exit 0.

---

## 4. Project Structure

New code follows the existing layout. Additions marked **+**.

```
cmd/
  web/                    → site + admin server
  admin/                  → create-admin CLI

pkg/
  handlers/               → Echo handlers; one struct per area, self-registering via init()
  + handlers/stickers.go      browse packs, pack detail, sticker JSON
  + handlers/upload.go        invite-gated submission form + POST
  + handlers/moderation.go    admin queue: list, approve, reject
  postgres/
    migrations/           → golang-migrate pairs; sqlc reads these as its schema
    queries/              → hand-written SQL; sqlc input
    + queries/packs.sql, stickers.sql, tags.sql, submissions.sql
    db/                   → GENERATED by sqlc; never hand-edit
  services/
  + services/storage.go       R2 client: Put, Copy, Delete, PublicURL, PresignedGet
  + services/zip.go           safe zip extraction + validation
  + services/imagemeta.go     format sniffing + animation detection
  ui/
    components/           → reusable gomponents nodes
    forms/                → form definitions (see pkg/ui/forms/file.go for the pattern)
    layouts/, icons/, models/, cache/
    pages/                → one exported func per page
  + ui/pages/stickers.go, upload.go, moderation.go
  + ui/components/sticker.go  sticker tile, pack card, copy button, favorite toggle
  routenames/             → route name constants; add new ones here

public/static/
  + js/db.js              IndexedDB wrapper (framework-free)
  + js/clipboard.js       copy dispatch: PNG-to-clipboard vs link (framework-free)
  + js/app.js             Alpine.data()/Alpine.store() registrations, inside alpine:init

docs/
  intent/meme-app.md      confirmed intent
  spec/phase-1.md         this file
tasks/
  plan.md, todo.md        produced by the next skill
```

---

## 5. Code Style

Match the surrounding code exactly. The dot-import of gomponents is the house style here, not
an accident — follow it.

**Views** (`pkg/ui/...`) — pure functions returning `Node`:

```go
package components

import (
	"github.com/tung-dnt/meme-app/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type StickerTileParams struct {
	ID       int64
	Name     string
	URL      string
	Animated bool
}

// StickerTile renders one sticker in a pack grid. The copy behaviour is decided
// client-side from data-animated, because an animated image cannot be written to
// the clipboard as an image (see spec §6.3).
func StickerTile(p StickerTileParams) Node {
	return Button(
		Class("sticker-tile"),
		Data("sticker-id", fmt.Sprint(p.ID)),
		Data("sticker-url", p.URL),
		Data("animated", strconv.FormatBool(p.Animated)),
		Img(Src(p.URL), Alt(p.Name), Loading("lazy")),
	)
}
```

**Handlers** (`pkg/handlers/...`) — struct + `init()` registration + `Init` + `Routes`, exactly as
`pkg/handlers/files.go` does:

```go
type Stickers struct {
	queries *pgdb.Queries
	storage *services.Storage
}

func init() { Register(new(Stickers)) }

func (h *Stickers) Init(c *services.Container) error {
	h.queries = c.Queries
	h.storage = c.Storage
	return nil
}

func (h *Stickers) Routes(g *echo.Group) {
	g.GET("/packs/:slug", h.Pack).Name = routenames.Pack
}
```

**Conventions**
- Route names always go in `pkg/routenames`; never inline a path string in a view.
- SQL lives in `pkg/postgres/queries/*.sql`; `pkg/postgres/db` is generated — never hand-edited.
- Errors returned up to Echo's handler; `msg.Success` / `msg.Error` for user-facing flash.
- Comments explain *why*, not *what* (the existing Makefile comments are the standard to match).
- Client JS: ES modules, no bundler, `const`/`let`, no globals except one namespace.
- Alpine components are registered in `app.js`, never defined inline in an `x-data`
  attribute — a view should read `x-data="stickerGrid"`, not carry a JS object literal.

**Alpine registration** (`public/static/js/app.js`) — note the `alpine:init` wrapper:

```js
import * as db from './db.js'
import { copySticker } from './clipboard.js'

document.addEventListener('alpine:init', () => {
  // Global favorite state, shared by the grid, the pack page and the favorites page.
  Alpine.store('favorites', {
    ids: new Set(),
    async load() { this.ids = new Set(await db.favoriteIds()) },
    has(id) { return this.ids.has(id) },
  })

  Alpine.data('stickerGrid', () => ({
    // One handler on the container, not one per tile: a pack can hold hundreds.
    async onClick(e) {
      const tile = e.target.closest('[data-sticker-id]')
      if (tile) await copySticker(tile.dataset)
    },
  }))
})
```

**Alpine in a view** — directives are plain attributes, so gomponents needs no helper:

```go
Div(
    Attr("x-data", "stickerGrid"),
    Attr("@click", "onClick($event)"),
    Class("sticker-grid"),
    Group(tiles),
)
```

---

## 6. Functional Requirements

### 6.1 Data model

New migration `000002_stickers`:

```sql
CREATE TABLE packs (
    id         BIGSERIAL PRIMARY KEY,
    slug       TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    is_global  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- The single global pack that every single-image upload lands in.
INSERT INTO packs (slug, name, is_global) VALUES ('singles', 'Singles', TRUE);

-- Guarantees there is never more than one global pack. A partial unique index is
-- the cheapest way to express "at most one row where is_global" in Postgres.
CREATE UNIQUE INDEX packs_one_global_idx ON packs (is_global) WHERE is_global;

CREATE TABLE stickers (
    id           BIGSERIAL PRIMARY KEY,
    pack_id      BIGINT NOT NULL REFERENCES packs (id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'pending',  -- pending | approved | rejected
    reject_note  TEXT,
    object_key   TEXT NOT NULL,       -- key within the R2 bucket
    mime_type    TEXT NOT NULL,
    animated     BOOLEAN NOT NULL,
    width        INT NOT NULL,
    height       INT NOT NULL,
    bytes        BIGINT NOT NULL,
    position     INT NOT NULL DEFAULT 0,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at  TIMESTAMPTZ
);

CREATE TABLE tags (
    id   BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE pack_tags (
    pack_id BIGINT NOT NULL REFERENCES packs (id) ON DELETE CASCADE,
    tag_id  BIGINT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (pack_id, tag_id)
);

CREATE INDEX stickers_status_idx   ON stickers (status, submitted_at);
CREATE INDEX stickers_pack_id_idx  ON stickers (pack_id, position);
CREATE INDEX stickers_approved_idx ON stickers (pack_id) WHERE status = 'approved';
CREATE INDEX pack_tags_tag_id_idx  ON pack_tags (tag_id);
```

**Approval status lives on `stickers`, not `packs`.** A pack is a container with no state of its
own; it becomes publicly visible once it has at least one approved sticker. This lets the admin
reject 2 bad images out of a 20-image zip instead of nuking the whole pack, and it makes the global
"Singles" pack — permanently a mix of pending, approved, and rejected content — an ordinary case
rather than a special one.

**Submissions map to packs like this:** a zip creates one new pack; a single image is appended to
the one global pack (`is_global = TRUE`, slug `singles`). Still one code path for moderation,
favoriting, and display.

**Approved-sticker counts are computed, not stored.** A denormalized counter would need updating on
every approve, reject, and delete; at this library size a `COUNT(*) FILTER (WHERE status =
'approved')` is cheaper than the bugs.

**No favorites table, no device-identity table, no users beyond the existing admin.** Favorites are
client-side only in phase 1. The existing `users` / `password_tokens` tables are retained unchanged
and hold exactly one row: the admin.

### 6.2 Browse, search, filter

- `GET /` — grid of pack cards, each with a cover sticker and its approved-sticker count. A pack
  appears only if it has ≥1 approved sticker. The global "Singles" pack renders as an ordinary
  card. Newest-approval-first, paginated with the existing `pkg/pager`.
- `GET /packs/:slug` — pack detail; grid of that pack's **approved** stickers, paginated — the
  global pack is the one guaranteed to outgrow a single page.
- `GET /search?q=&tag=` — `ILIKE` over `packs.name`, `stickers.name`, `tags.name`, restricted to
  `stickers.status = 'approved'`. A `pg_trgm` GIN index is **deliberately deferred** until the
  library is large enough for sequential scans to hurt; noted in Open Questions.
- Tag list rendered as filter chips.
- Only `stickers.status = 'approved'` is ever visible to a visitor, enforced in the SQL rather than
  in handler code so that no future caller can forget it. Pack visibility is *derived* from that
  join and never stored, so the two can't drift out of sync.

### 6.3 Copy to clipboard — the two paths

This is the core interaction and it has a hard platform limit.

| Sticker | Mechanism | Result |
|---|---|---|
| Static (PNG/JPEG/static WebP) | `navigator.clipboard.write([new ClipboardItem({'image/png': blob})])` | Image on clipboard; pastes as an image |
| Animated (GIF/APNG/animated WebP) | `navigator.clipboard.writeText(absoluteURL)` | URL on clipboard; chat apps unfurl it, animation preserved |

Requirements:
- `animated` is determined **server-side at upload** and stored on the row — never guessed by the
  client. Detection: GIF → more than one image descriptor; PNG → presence of an `acTL` chunk;
  WebP → presence of an `ANIM` chunk.
- Non-PNG static images are converted to PNG in-browser via `canvas` before the clipboard write,
  because `image/png` is the only type with universal `ClipboardItem` support.
- **Safari constraint:** `ClipboardItem` must be constructed with a *Promise* for the blob,
  synchronously inside the click handler's task. Awaiting the fetch first and *then* calling
  `clipboard.write()` loses the user-gesture context and Safari silently rejects it. The
  implementation must pass `new ClipboardItem({'image/png': fetchAndConvert()})` without awaiting.
- The UI must show which path ran — a toast reading "Copied image" vs "Copied link" — so the
  behaviour is never a silent surprise.
- Clipboard access requires a secure context. HTTPS in production is therefore a hard requirement,
  not a nice-to-have.

### 6.4 Favorites — local only

- IndexedDB database `memeapp`, stores: `favorites` (keyed by sticker id), `favorite_packs`
  (keyed by pack id), `blobs` (keyed by object key).
- Favoriting a sticker stores its metadata **and fetches + caches the image blob**, so the
  favorites page renders from local data with no network round-trip.
- Favoriting a pack stores the pack and all of its **approved** stickers' blobs.
- **The global "Singles" pack cannot be favorited as a pack** — its stickers are favoritable
  individually only. It grows without bound, so a one-tap "cache everything in here" would be an
  unbounded download that breaks the blob budget below. The pack-favorite control is simply not
  rendered when `is_global`.
- `GET /favorites` returns an empty shell; the page is hydrated entirely by `favorites.js` from
  IndexedDB. The server never learns what anyone favorited.
- Cache eviction: hard cap on cached blob bytes (default 100 MB); evict least-recently-favorited
  first. Metadata is never evicted, so a favorite survives eviction and simply re-fetches.
- **Known and accepted limitation:** favorites are per-browser. Clearing site data loses them, and
  Safari's ITP may evict IndexedDB after 7 days without interaction. Cross-device sync via passkey
  is phase 2. The UI states this plainly on the favorites page rather than pretending otherwise.

### 6.5 Upload

- `GET /upload` — form: invite code, pack name, tags, file.
- `POST /upload` — accepts one file: an image, or a zip.
- **One zip = one pack.** Pack name comes from the single top-level folder inside the zip; if
  entries sit at the zip root, the zip's filename is used. An explicit name in the form overrides
  both.
- **A single image is appended to the global "Singles" pack**, not given a pack of its own. The
  pack-name field is therefore shown only for zip uploads.
- Validation, all enforced server-side before anything is written to R2:

  | Guard | Limit |
  |---|---|
  | Request body | 100 MB |
  | Zip entries | 200 |
  | Total uncompressed | 250 MB |
  | Per-image | 5 MB |
  | Compression ratio | reject above 100:1 |
  | Dimensions | 32×32 min, 2048×2048 max |
  | Accepted types | PNG, JPEG, GIF, WebP — by **content sniffing**, never by file extension |

- Path traversal: zip entry paths are discarded entirely; only `filepath.Base` is used. Directory
  entries, symlinks, and dotfiles are skipped.
- Any entry failing validation is skipped with a per-file reason surfaced to the uploader; the
  submission still succeeds if at least one image is valid.
- Zip extraction runs as a `backlite` task so the request returns promptly; the submission appears
  in the queue as pending stickers once extraction completes.
- **Invite code** lives in config (`app.inviteCode`, overridable by env). Compared with
  `subtle.ConstantTimeCompare`. Rate-limited to 10 attempts/hour/IP.

### 6.6 Moderation

- All routes under `/admin`, behind the existing `middleware.RequireAuthentication` plus an
  admin check.
- `GET /admin/queue` — pending **stickers**, oldest first, **grouped by pack**, with a preview grid.
- `POST /admin/stickers/:id/approve` — sets `status = 'approved'`, copies that object from the
  pending bucket to the public bucket, deletes the pending object.
- `POST /admin/stickers/:id/reject` — sets `status = 'rejected'` with an optional note, deletes the
  pending object.
- `POST /admin/packs/:id/approve` — bulk: applies the above to every pending sticker in the pack.
  This is the common case for a zip and must be one click, not twenty.
- Partial approval is a first-class outcome: rejecting 2 stickers from a 20-image zip leaves the
  other 18 publicly visible.
- Public registration is **removed**: the `/register`, `/forgot-password`, and `/reset-password`
  routes and their pages/forms are deleted. `/login` and `/logout` remain, for the admin only.

### 6.7 Object storage (R2)

Two buckets, because R2 public access is per-bucket and pending content must not be reachable:

| Bucket | Access | Contents |
|---|---|---|
| `stickers-pending` | private | `pending/{sticker_id}.{ext}` |
| `stickers-public` | public (custom domain) | `packs/{pack_id}/{sticker_id}.{ext}` |

- Objects move **per sticker**, on that sticker's approval — never per pack.
- Admin previews pending images via presigned GET URLs (15 min TTL).
- Approved images are served **directly from the public bucket's custom domain** — image bytes
  never pass through the Go app, which is the reason for choosing R2.
- `Cache-Control: public, max-age=31536000, immutable` on approved objects; keys are immutable, so
  a change means a new key.
- CORS on the public bucket must allow `GET` from the site origin — `fetch()` for the clipboard
  and blob caching needs it, and a plain `<img>` load succeeding does **not** imply `fetch` will.
- New config block `storage:` with account ID, bucket names, access key, secret, and public base
  URL. Credentials come from env, never `config.yaml`.

---

## 7. Testing Strategy

Match the existing setup: Go stdlib `testing`, helpers in `pkg/tests`, real PostgreSQL against the
`$RAND` throwaway schema, run via `make test`.

| Level | Location | Covers |
|---|---|---|
| Unit | `pkg/services/*_test.go` | zip validation, animation detection, format sniffing, slug generation |
| Query | `pkg/postgres/db` via `pkg/tests` | every sqlc query, especially that visitor queries cannot return non-approved stickers or empty packs |
| Handler | `pkg/handlers/*_test.go` | routing, invite-code gate, admin auth gate, approve/reject transitions |
| Manual | browser | clipboard (both paths), IndexedDB, offline favorites |

**Required negative tests** — these encode the guards that matter:
- A zip bomb (high ratio) is rejected.
- A zip with `../../etc/passwd` entries writes nothing outside the target.
- A `.png` file whose bytes are actually a PHP script is rejected by sniffing.
- An anonymous request to any `/admin` route is redirected, not served.
- A pending sticker appears in none of browse, search, or pack-detail; a pack whose stickers are
  all pending does not appear in the grid at all.
- A second row with `is_global = TRUE` is rejected by the database.

**Client-side JS has no automated tests in phase 1.** This is a deliberate, time-driven cut and the
single largest quality risk in the plan. It is mitigated by a written manual checklist covering
Chrome, Firefox, Safari desktop, and iOS Safari — the last of which is where clipboard behaviour
most often differs.

---

## 8. Boundaries

**Always**
- Regenerate with `make sqlc-gen` after touching any `.sql`; commit the generated output.
- Write migrations as reversible `.up.sql`/`.down.sql` pairs.
- Validate uploads by content sniffing, never by file extension.
- Enforce `stickers.status = 'approved'` in SQL, not in handler code.
- Add route names to `pkg/routenames`.
- Run `make test` before committing.

**Ask first**
- Adding any dependency beyond `aws-sdk-go-v2` (already approved for R2).
- Changing the `users` / `password_tokens` schema.
- Anything that puts image bytes through the Go app instead of the public bucket.
- Introducing a JS build step or a frontend framework.
- Changing the two-bucket topology.

**Never**
- Commit R2 credentials, the invite code, or the encryption key.
- Hand-edit `pkg/postgres/db/` (generated).
- Serve or link a sticker that is not `approved`.
- Render a pack-favorite control on the global pack.
- Reintroduce device fingerprinting — it cannot work across devices and collides on a public site
  (see `docs/intent/meme-app.md`, decision 1).
- Ship a copy button that silently no-ops on animated stickers.

---

## 9. Success Criteria

Phase 1 is done when all of the following are demonstrably true:

1. A visitor with no cookie or session can browse the pack grid and open a pack detail page,
   seeing only approved stickers.
2. Clicking a **static** sticker puts a PNG on the clipboard that pastes as an image into Slack,
   Discord, and iOS Messages.
3. Clicking an **animated** sticker puts its URL on the clipboard, and pasting it into the same
   apps yields a moving image.
4. The UI states which of the two happened, every time.
5. Favoriting a sticker, then going fully offline, still renders that sticker on `/favorites`.
6. Favoriting a pack caches all of its stickers.
7. Search by name and filter by tag each return correct results over approved content only.
8. A contributor with the invite code can upload a zip of 20 images; it appears in the admin queue
   as one group of 20 pending stickers, approvable in a single click.
9. A single-image upload lands in the global "Singles" pack — no new pack is created — and appears
   there once approved.
10. Rejecting 2 stickers out of a 20-image zip leaves the other 18 publicly visible.
11. The "Singles" pack renders no favorite-pack control.
12. A contributor without the code cannot upload; the attempt is rate-limited.
13. Approving a sticker moves its object to the public bucket; rejecting deletes it. A pack becomes
    visible exactly when its first sticker is approved.
14. Every negative test in §7 passes.
15. `make sqlc-vet && make test && make build` exits 0, and CI is green.
16. The site is deployed over HTTPS with a real domain.

---

## 10. Open Questions

1. **Go version mismatch.** `go.mod` declares `go 1.27.0`; `.github/workflows/test.yml` pins
   `go-version: '1.25'`. CI is testing a different toolchain than the one the module targets.
   Bump CI to 1.27 as part of phase 1?
2. **`pg_trgm`.** Deferred by default. Add the extension + GIN index now, or wait for a measured
   slow query?
3. **Invite code rotation.** Config-only means rotation requires a redeploy. Acceptable, or should
   it live in a DB table the admin can edit?
4. **Public bucket domain.** Custom domain on R2, or the default `r2.dev` URL for phase 1?
5. **Rejected stickers.** Retain the DB row for audit, or hard-delete after N days?
6. **iOS Safari clipboard.** Needs early hands-on verification — it is the most likely place the
   copy interaction breaks, and it would be expensive to discover late. Worth spiking before the
   rest of the UI is built?
