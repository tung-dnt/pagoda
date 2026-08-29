# T26 · Bulk pack approve

| | |
|---|---|
| **Epic** | [E07 · Moderation](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T25 |
| **Unlocks** | T28 |
| **Spec** | §6.6 |
| **Status** | ☐ not started |

## Problem

`POST /admin/packs/:id/approve` applies T25's approve to every pending sticker in a pack.
Approving a 20-image zip must be one click, not twenty.

## Given

- Single-sticker approve from T25

## Constraints

- Reuse T25's per-sticker logic. Do not write a second promote path that can drift from the
  first.
- **Partial failure is expected.** If 18 promote and 2 fail, report exactly that and leave
  the 18 approved. All-or-nothing would mean one flaky network call discards good work.
- Never applies to the global "Singles" pack — its contents come from unrelated
  contributors and reviewing them as a batch defeats the point. Hide the control there.
- Bulk-approving 200 stickers is 200 storage round-trips. Bound the concurrency (say 5) so
  one click cannot saturate the connection pool.
- Partial approval stays first-class: rejecting 2 stickers first, then bulk-approving,
  leaves 18 approved and 2 rejected.

## Acceptance

- [ ] One click approves every pending sticker in a pack
- [ ] The pack becomes visible with the correct count
- [ ] A partial failure reports which stickers failed and keeps the successes
- [ ] The control is absent on the "Singles" pack
- [ ] Rejecting 2 of 20 first, then bulk-approving, yields 18 approved and 2 rejected
- [ ] Concurrency is bounded

## Verify

```
make test && make run
```

Upload a 20-image zip → reject 2 in the queue → bulk approve → confirm `/packs/:slug` shows
18 and the pack card count reads 18.

## Files

- `pkg/handlers/moderation.go`, `pkg/handlers/moderation_test.go`
- `pkg/ui/pages/moderation.go`

## Hints

- An `errgroup` with `SetLimit(5)` gives bounded concurrency and per-item errors in a few
  lines.
- This closes success criterion 10 in the spec — the partial-approval case. Test it
  explicitly.
