# T23 · Async zip ingestion task

| | |
|---|---|
| **Epic** | [E06 · Upload](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T22, T20 |
| **Unlocks** | T24 |
| **Spec** | §6.5 |
| **Status** | ☐ not started |

## Problem

Handle `POST /upload` for a **zip**: accept it, return promptly, and process it in a
background `backlite` task that creates one new pack with N pending stickers.

Extraction of a 200-image archive must not happen inside the request.

## Given

- `Extract` from T20, upload handler from T22
- `pkg/tasks/example.go` + `register.go` — the backlite task pattern (**delete the example
  once this works**, per T03)
- `c.Tasks` on the container

## Constraints

- The request stores the raw zip somewhere the task can reach it, then returns. Use the
  pending bucket or local `afero` scratch — either is fine, but **clean it up** on both
  success and failure.
- The task is retried by backlite on failure, so it must be **idempotent**. A retry after a
  partial run must not create a second pack or duplicate stickers. Key it on the submission
  id, not on a name.
- One zip = one new pack. Name from the top-level folder, or the zip filename, overridden by
  the form field.
- Per-entry skip reasons must reach the uploader. If they only ever land in a log, a
  contributor whose 3 images were silently dropped has no way to know.
- A zip with zero valid images fails the submission with a clear reason and creates no pack.
- Tags from the form attach to the new pack, slug-normalized the same way as T17.

## Acceptance

- [ ] A 20-image zip returns promptly and yields one pack with 20 pending stickers
- [ ] Re-running the task creates no duplicates
- [ ] Skip reasons are visible to the uploader
- [ ] A zero-valid-image zip creates no pack and reports why
- [ ] The temporary zip is deleted on both paths
- [ ] The new pack is invisible in browse until a sticker is approved
- [ ] Tags are attached to the new pack

## Verify

```
make test && make run
```

Upload a 20-image zip → response is immediate → watch `/admin/tasks` (the backlite UI is
already wired) → confirm 20 pending rows and one invisible pack. Then re-run the task from
that UI and confirm no duplicates.

## Files

- `pkg/tasks/ingest.go`, `pkg/tasks/register.go`
- `pkg/handlers/upload.go`
- `pkg/postgres/queries/packs.sql`

## Hints

- The existing `/admin/tasks` backlite UI is genuinely useful here — you get retries and
  failure inspection for free.
- Idempotency: insert the pack with the submission id as a unique key, and let a duplicate
  insert fail harmlessly.
