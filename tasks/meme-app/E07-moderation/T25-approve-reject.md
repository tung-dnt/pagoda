# T25 · Approve / reject a sticker

| | |
|---|---|
| **Epic** | [E07 · Moderation](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T24 |
| **Unlocks** | T26 |
| **Spec** | §6.6 |
| **Status** | ☐ not started |

## Problem

Wire the two actions. Approving promotes the object from the pending bucket to the public
one and flips the row to `approved`. Rejecting deletes the object and records why.

This is the task that makes a pack appear on the public site, so the ordering of DB and
storage operations matters.

## Given

- Queue page from T24, `Promote`/`Delete` from T05

## Constraints

- **Promote the object first, then flip the row.** If the row flips first and the copy
  fails, the site advertises a sticker whose bytes are not public — a visible broken image.
  The reverse failure is invisible and harmless.
- Both actions are `POST`. A `GET` that mutates will be fired by a crawler or a prefetcher.
- CSRF applies. It is already global.
- Reject records an optional note and sets `reviewed_at`. Keep the row — audit trail. (Spec
  Open Question 5 asks whether to purge later; not this task.)
- Approving a sticker in a pack with no previously-approved stickers makes that pack appear
  in browse. Confirm that actually happens — it is the derived-visibility rule from T04
  paying off, and the place it would silently fail.
- Actions must be idempotent. A double-submitted approve must not promote twice or error.
- Use HTMX to swap the row out in place; a full page reload after each of 20 decisions is
  miserable.
- If a swapped fragment carries `x-data`, Alpine 3 initializes it automatically — but state
  in the *replaced* node is discarded. Keep any state that must outlive a swap (an open
  reject-note box, a filter) on an ancestor outside the swap target, or it resets on every
  decision.

## Acceptance

- [ ] Approve promotes the object, sets `approved`, sets `reviewed_at`
- [ ] The approved sticker is immediately visible in browse and its public URL works
- [ ] Reject deletes the pending object and sets `rejected` with the note
- [ ] A rejected sticker never appears publicly and its object is gone
- [ ] Approving the first sticker in a pack makes that pack appear on `/`
- [ ] Double-approve is harmless
- [ ] `GET` on either action route is not accepted

## Verify

```
make test && make run
```

Approve one sticker → check the public URL with `curl -I` → confirm 200 and the
`Cache-Control: immutable` header → confirm the pack now appears on `/`. Reject another →
confirm the pending object is gone and it is absent everywhere.

## Files

- `pkg/handlers/moderation.go`, `pkg/handlers/moderation_test.go`
- `pkg/ui/pages/moderation.go`
- `pkg/postgres/queries/stickers.sql`

## Hints

- Do the promote outside the transaction, then the row update inside it. Storage cannot
  participate in a DB transaction, so pick the failure mode you can live with.
- HTMX: `hx-post` on the button with `hx-target="closest [data-sticker-row]"` and
  `hx-swap="outerHTML"`.
