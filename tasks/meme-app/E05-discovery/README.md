# E05 · Search & tags

Find a specific sticker by name or tag.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T17](./T17-tag-queries.md) | Tag queries + chips | S | ~1.5 h | T04, T08 |
| [T18](./T18-search.md) | Search endpoint + results | M | ~2.5 h | T17 |
| [T19](./T19-tag-filter.md) | Tag filtering | S | ~1.5 h | T18 |

## Checkpoint — Search & tags

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
