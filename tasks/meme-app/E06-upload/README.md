# E06 · Upload

Invite-gated submission of a single image or a zip.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T20](./T20-zip-extract.md) | Safe zip extraction service | M | ~3 h | T06 |
| [T21](./T21-invite-gate.md) | Upload form + invite gate + rate limit | M | ~2.5 h | T05 |
| [T22](./T22-upload-single.md) | Upload POST: single image → Singles | M | ~2.5 h | T21, T06 |
| [T23](./T23-upload-zip.md) | Async zip ingestion task | M | ~2.5 h | T22, T20 |

## Checkpoint — Upload

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
