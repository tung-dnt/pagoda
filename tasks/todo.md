# TODO: Meme Sticker App — Phase 1

29 tasks. One task ≈ one sitting. Full detail in [`tasks/meme-app/`](./meme-app/);
plan and risks in [`tasks/plan.md`](./plan.md).

**Sizes:** XS ~30 min · S ~1 h · M ~2–3 h. Nothing larger is allowed.

**Rule:** every task ends with `make test && make build` exiting 0.


## E01 · Foundation

- [ ] **[T01](./meme-app/E01-foundation/T01-clipboard-spike.md)** · Clipboard spike on iOS Safari — `S` · deps: None
- [ ] **[T02](./meme-app/E01-foundation/T02-strip-public-auth.md)** · Strip public registration — `S` · deps: None
- [ ] **[T03](./meme-app/E01-foundation/T03-strip-demo-pages.md)** · Strip boilerplate demo pages — `S` · deps: T02
- [ ] **[T04](./meme-app/E01-foundation/T04-schema.md)** · Schema migration + sqlc queries — `M` · deps: None
- [ ] **[T05](./meme-app/E01-foundation/T05-r2-storage.md)** · R2 storage service — `M` · deps: None
- [ ] **Checkpoint** — [E01](./meme-app/E01-foundation/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E02 · Browse

- [ ] **[T06](./meme-app/E02-browse/T06-image-metadata.md)** · Image sniffing + animation detection — `S` · deps: None
- [ ] **[T07](./meme-app/E02-browse/T07-seed.md)** · Dev seed command — `S` · deps: T04, T05, T06
- [ ] **[T08](./meme-app/E02-browse/T08-home-grid.md)** · Home page: pack grid — `M` · deps: T04, T07
- [ ] **[T09](./meme-app/E02-browse/T09-pack-detail.md)** · Pack detail page — `M` · deps: T08
- [ ] **Checkpoint** — [E02](./meme-app/E02-browse/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E03 · Copy to clipboard

- [ ] **[T10](./meme-app/E03-clipboard/T10-clipboard-js.md)** · clipboard.js: the two copy paths — `M` · deps: T01, T09
- [ ] **[T11](./meme-app/E03-clipboard/T11-copy-wiring.md)** · Wire copy into the sticker tile — `S` · deps: T10
- [ ] **Checkpoint** — [E03](./meme-app/E03-clipboard/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E04 · Favorites

- [ ] **[T12](./meme-app/E04-favorites/T12-indexeddb.md)** · db.js: IndexedDB wrapper — `M` · deps: None
- [ ] **[T13](./meme-app/E04-favorites/T13-favorite-toggle.md)** · favorites.js: toggle + blob cache — `M` · deps: T12, T11
- [ ] **[T14](./meme-app/E04-favorites/T14-favorites-page.md)** · /favorites page + hydration — `M` · deps: T13
- [ ] **[T15](./meme-app/E04-favorites/T15-favorite-pack.md)** · Pack-level favorite — `S` · deps: T14
- [ ] **[T16](./meme-app/E04-favorites/T16-service-worker.md)** · Service worker: offline shell — `M` · deps: T14, T15
- [ ] **Checkpoint** — [E04](./meme-app/E04-favorites/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E05 · Search & tags

- [ ] **[T17](./meme-app/E05-discovery/T17-tag-queries.md)** · Tag queries + chips — `S` · deps: T04, T08
- [ ] **[T18](./meme-app/E05-discovery/T18-search.md)** · Search endpoint + results — `M` · deps: T17
- [ ] **[T19](./meme-app/E05-discovery/T19-tag-filter.md)** · Tag filtering — `S` · deps: T18
- [ ] **Checkpoint** — [E05](./meme-app/E05-discovery/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E06 · Upload

- [ ] **[T20](./meme-app/E06-upload/T20-zip-extract.md)** · Safe zip extraction service — `M` · deps: T06
- [ ] **[T21](./meme-app/E06-upload/T21-invite-gate.md)** · Upload form + invite gate + rate limit — `M` · deps: T05
- [ ] **[T22](./meme-app/E06-upload/T22-upload-single.md)** · Upload POST: single image → Singles — `M` · deps: T21, T06
- [ ] **[T23](./meme-app/E06-upload/T23-upload-zip.md)** · Async zip ingestion task — `M` · deps: T22, T20
- [ ] **Checkpoint** — [E06](./meme-app/E06-upload/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E07 · Moderation

- [ ] **[T24](./meme-app/E07-moderation/T24-queue.md)** · Admin queue page — `M` · deps: T23
- [ ] **[T25](./meme-app/E07-moderation/T25-approve-reject.md)** · Approve / reject a sticker — `M` · deps: T24
- [ ] **[T26](./meme-app/E07-moderation/T26-bulk-approve.md)** · Bulk pack approve — `S` · deps: T25
- [ ] **Checkpoint** — [E07](./meme-app/E07-moderation/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## E08 · Ship

- [ ] **[T27](./meme-app/E08-ship/T27-ci.md)** · Fix CI Go version + green pipeline — `XS` · deps: None
- [ ] **[T28](./meme-app/E08-ship/T28-browser-checklist.md)** · Manual browser verification — `S` · deps: T26
- [ ] **[T29](./meme-app/E08-ship/T29-deploy.md)** · Deploy: HTTPS, buckets, CORS, env — `M` · deps: T27, T28
- [ ] **Checkpoint** — [E08](./meme-app/E08-ship/README.md): all tasks done, `make sqlc-vet && make test && make build` green

## Done means

All 16 success criteria in [spec §9](../docs/spec/phase-1.md) verified against production,
not against localhost.
