# T03 · Strip boilerplate demo pages

| | |
|---|---|
| **Epic** | [E01 · Foundation](./README.md) |
| **Size** | S — ~1 h |
| **Depends on** | T02 |
| **Unlocks** | T08 |
| **Spec** | — |
| **Status** | ☐ not started |

## Problem

Pagoda ships demo pages that have nothing to do with this app: About, Contact, Cache,
Task, Files, and a **fake** search handler that returns lorem ipsum. Leaving them means
every future grep for "search" hits a decoy, and the home page is someone else's.

Remove them now, while the repo is still small enough that it's a 1-hour job.

## Given

- `pkg/handlers/`: `contact.go`, `cache.go`, `task.go`, `files.go`, `search.go`, `pages.go`
- `pkg/ui/pages/`: `about.go`, `contact.go`, `cache.go`, `task.go`, `file.go`, `search.go`, `home.go`
- `pkg/ui/forms/`: `contact.go`, `cache.go`, `task.go`, `file.go`
- `pkg/ui/models/post.go` — fake blog posts used by the demo home page

## Constraints

- **Keep `pkg/handlers/files.go` as a reference until T22 is done**, then delete it. It is
  the only working example of multipart upload with `afero` in the repo. Copy what you need
  out of it first. Note that decision in the commit message.
- Keep `pkg/handlers/error.go` and `pkg/handlers/pages.go`'s home handler — home gets
  rewritten in T08, not deleted, so the route name `Home` survives.
- Keep `pkg/ui/components/` wholesale. Those are generic and T08 will use them.
- Keep the `Search` route name; T18 reuses it for the real thing.
- `pkg/tasks/example.go` is a backlite demo — keep it until T23 shows the real task pattern.

## Acceptance

- [ ] `/about`, `/contact`, `/cache`, `/task` return 404
- [ ] `/search` no longer returns lorem ipsum (a stub or 404 is fine for now)
- [ ] `pkg/ui/models/post.go` is gone and nothing references it
- [ ] Nav contains no links to removed pages
- [ ] `make test && make build` is green

## Verify

```
make test && make build
grep -rn "lorem\|Lorem" pkg/ | grep -v _test.go   # expect no hits
```

## Files

- roughly 12 deletions across `pkg/handlers/`, `pkg/ui/pages/`, `pkg/ui/forms/`
- `pkg/handlers/router_test.go`, `pkg/handlers/pages_test.go`

## Hints

- Handlers self-register via `init()` + `Register(new(X))`, so deleting the file is enough
  to unregister the routes. There is no central list to edit.
- Do this in one commit so the diff is reviewable as "delete demo content".
