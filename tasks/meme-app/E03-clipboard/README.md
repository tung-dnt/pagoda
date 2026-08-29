# E03 · Copy to clipboard

The core interaction. Two paths: PNG for static, link for animated.

| Task | Title | Size | Budget | Depends on |
|---|---|---|---|---|
| [T10](./T10-clipboard-js.md) | clipboard.js: the two copy paths | M | ~2.5 h | T01, T09 |
| [T11](./T11-copy-wiring.md) | Wire copy into the sticker tile | S | ~1.5 h | T10 |

## Checkpoint — Copy to clipboard

Do not start the next epic until all of these hold:

- [ ] Every task above is checked off
- [ ] `make sqlc-vet && make test && make build` exits 0
- [ ] No task left a TODO that the next epic silently depends on

[← back to the board](../README.md)
