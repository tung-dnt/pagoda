# E01 · Foundation

De-risk the riskiest unknown, strip what we don't need, and lay the schema + storage everything else sits on.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T01](./T01-clipboard-spike.md) | Clipboard spike on iOS Safari | S | ~1 h | None |
| [T02](./T02-strip-public-auth.md) | Strip public registration | S | ~1 h | None |
| [T03](./T03-strip-demo-pages.md) | Strip boilerplate demo pages | S | ~1 h | T02 |
| [T04](./T04-schema.md) | Schema migration + sqlc queries | M | ~3 h | None |
| [T05](./T05-r2-storage.md) | R2 storage service | M | ~3 h | None |

## Checkpoint — Foundation

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
