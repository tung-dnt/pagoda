# E04 · Favorites

Local-first favorites in IndexedDB, working offline.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T12](./T12-indexeddb.md) | db.js: IndexedDB wrapper | M | ~2.5 h | None |
| [T13](./T13-favorite-toggle.md) | favorites.js: toggle + blob cache | M | ~2.5 h | T12, T11 |
| [T14](./T14-favorites-page.md) | /favorites page + hydration | M | ~2.5 h | T13 |
| [T15](./T15-favorite-pack.md) | Pack-level favorite | S | ~1.5 h | T14 |
| [T16](./T16-service-worker.md) | Service worker: offline shell | M | ~3 h | T14, T15 |

## Checkpoint — Favorites

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
