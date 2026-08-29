# E08 · Ship

Green CI, verified in real browsers, deployed over HTTPS.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T27](./T27-ci.md) | Fix CI Go version + green pipeline | XS | ~30 min | None |
| [T28](./T28-browser-checklist.md) | Manual browser verification | S | ~1.5 h | T26 |
| [T29](./T29-deploy.md) | Deploy: HTTPS, buckets, CORS, env | M | ~3 h | T27, T28 |

## Checkpoint — Ship

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
