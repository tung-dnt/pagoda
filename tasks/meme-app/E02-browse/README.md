# E02 · Browse

First visible slice: real stickers rendering in a browser.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T06](./T06-image-metadata.md) | Image sniffing + animation detection | S | ~1.5 h | None |
| [T07](./T07-seed.md) | Dev seed command | S | ~1.5 h | T04, T05, T06 |
| [T08](./T08-home-grid.md) | Home page: pack grid | M | ~2.5 h | T04, T07 |
| [T09](./T09-pack-detail.md) | Pack detail page | M | ~2.5 h | T08 |

## Checkpoint — Browse

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
