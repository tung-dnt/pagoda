# E07 · Moderation

Per-sticker approve/reject with one-click bulk approve for a zip.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T24](./T24-queue.md) | Admin queue page | M | ~2.5 h | T23 |
| [T25](./T25-approve-reject.md) | Approve / reject a sticker | M | ~2.5 h | T24 |
| [T26](./T26-bulk-approve.md) | Bulk pack approve | S | ~1.5 h | T25 |

## Checkpoint — Moderation

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
