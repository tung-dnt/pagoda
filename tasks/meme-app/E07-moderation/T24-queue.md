# T24 · Admin queue page

| | |
|---|---|
| **Epic** | [E07 · Moderation](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T23 |
| **Unlocks** | T25 |
| **Spec** | §6.6 |
| **Status** | ☐ not started |

## Problem

`GET /admin/queue`: every pending sticker, oldest first, grouped by pack, with a visual
preview of each. This is where the admin decides what goes public, so they must be able to
actually *see* the images — pending objects live in a private bucket.

## Given

- `pkg/handlers/admin.go` — existing admin group with `middleware.RequireAdmin`
- `PresignedGet` from T05
- `pkg/pager`

## Constraints

- Everything under `/admin` is already behind `RequireAdmin`. Verify it with a test rather
  than trusting it.
- Previews use **presigned URLs, 15 min TTL**. Do not make the pending bucket public, and
  do not proxy bytes through the app.
- Presigning N URLs per page render is N signature computations — keep the page size
  modest (25) and generate them lazily where you can.
- Group by pack so a 20-image zip reads as one unit, not 20 unrelated rows.
- Show what the admin needs to judge: image, name, dimensions, size, animated flag,
  submission time.
- No approve/reject actions yet. Display only. That is T25.

## Acceptance

- [ ] `/admin/queue` lists pending stickers grouped by pack, oldest first
- [ ] Every preview image actually renders
- [ ] An anonymous request is redirected to login, not served
- [ ] A logged-in non-admin (if one existed) is rejected
- [ ] Paginated at 25
- [ ] Approved and rejected stickers do not appear

## Verify

```
make test && make run
```

Log in as admin, open `/admin/queue`, confirm images render. Log out, hit the URL directly,
confirm redirect. Wait 16 minutes and reload — expired URLs should regenerate, not 403.

## Files

- `pkg/handlers/moderation.go`, `pkg/handlers/moderation_test.go`
- `pkg/ui/pages/moderation.go`
- `pkg/postgres/queries/stickers.sql`
- `pkg/routenames/names.go`

## Hints

- `pkg/handlers/admin.go` already shows the group pattern:
  `ag := g.Group("/admin", middleware.RequireAdmin)`.
- The anonymous-access test is cheap and catches the worst possible bug in this app. Write
  it even though the middleware "obviously" works.
