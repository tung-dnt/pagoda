# Feature: Meme Sticker App — Phase 1

Day-by-day board. One task ≈ one sitting. Work top to bottom; a task is only unblocked
when everything in its **Depends on** list is checked.

- Intent: [`docs/intent/meme-app.md`](../../docs/intent/meme-app.md)
- Spec: [`docs/spec/phase-1.md`](../../docs/spec/phase-1.md)
- Plan + risks: [`tasks/plan.md`](../plan.md)
- Flat checklist: [`tasks/todo.md`](../todo.md)

## Size legend

| Size | Budget | Meaning |
|---|---|---|
| **XS** | ~30 min | One config or one function |
| **S** | ~1 h | One component, one endpoint, one file |
| **M** | ~2–3 h | One evening. A whole slice |
| **L** | — | Not allowed. Split it |

## Epics

| # | Epic | Tasks | Budget | Delivers |
|---|---|---|---|---|
| E01 | [Foundation](./E01-foundation/) | 5 | ~9 h | Schema, storage, clean slate, clipboard proven |
| E02 | [Browse](./E02-browse/) | 4 | ~8 h | Stickers visible in a browser |
| E03 | [Copy to clipboard](./E03-clipboard/) | 2 | ~4 h | The core loop works |
| E04 | [Favorites](./E04-favorites/) | 5 | ~12 h | Offline favorites, incl. cold reload |
| E05 | [Search & tags](./E05-discovery/) | 3 | ~6 h | Findability |
| E06 | [Upload](./E06-upload/) | 4 | ~10 h | Contributors can submit |
| E07 | [Moderation](./E07-moderation/) | 3 | ~7 h | Nothing embarrassing goes public |
| E08 | [Ship](./E08-ship/) | 3 | ~6 h | Live on HTTPS |

**Total: 29 tasks, ~62 h.** At one evening a night that is roughly six weeks.

## Rule of the board

Every task ends with the repo **green** — `make test && make build` exits 0. If a task
can't end green, it was too big; split it and update this board.
