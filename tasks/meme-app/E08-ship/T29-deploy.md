# T29 · Deploy: HTTPS, buckets, CORS, env

| | |
|---|---|
| **Epic** | [E08 · Ship](./README.md) |
| **Size** | M — ~3 h |
| **Depends on** | T27, T28 |
| **Unlocks** | — |
| **Spec** | §6.7, §9 |
| **Status** | ☐ not started |

## Problem

Get it live: two R2 buckets configured correctly, secrets in the environment, the app
behind HTTPS on a real domain. Phase 1 is not done until someone who is not you can open it
on their phone.

## Given

- `services.Storage` from T05
- `docker-compose.yml`, `config/config.yaml`
- `make build`

## Constraints

- **HTTPS is mandatory, not a nice-to-have.** `navigator.clipboard` is undefined without a
  secure context, so the core interaction simply does not exist over plain HTTP.
- Two buckets: `stickers-pending` **private**, `stickers-public` **public** via a custom
  domain.
- **CORS on the public bucket** must allow `GET` from the site origin. An `<img>` that loads
  fine does *not* mean `fetch()` will succeed — and copy and blob caching both use `fetch`.
  This is the single most likely production-only failure. Test it deliberately.
- Every secret from the environment: R2 access key and secret, the invite code, and
  `app.encryptionKey`. **The default encryption key in `config.yaml` is committed and must
  be replaced** — sessions are signed with it.
- `app.host` must be the real public URL. The animated copy path builds absolute URLs from
  it, and a wrong value copies links that go nowhere.
- Set `app.environment` to something other than `local` — `make seed` refuses to run
  outside local, and that is the guard protecting production data.
- Confirm the trusted-proxy setup so `ctx.RealIP()` is honest, or T21's rate limiter keys
  everything to one address.
- Back up the database. It holds the only record of what has been approved.

## Acceptance

- [ ] The site is reachable over HTTPS on a real domain
- [ ] `navigator.clipboard` is defined in production
- [ ] Copy works, both paths, from a phone on mobile data
- [ ] Pending objects are **not** publicly reachable; approved ones are
- [ ] CORS allows `fetch` from the site origin — verified, not assumed
- [ ] No secret is present in any committed file
- [ ] The encryption key is not the committed default
- [ ] Migrations applied; `make seed` refuses to run
- [ ] Rate limiting keys on the real client IP
- [ ] A database backup exists and has been restored once as a test

## Verify

```
curl -I https://<domain>/
curl -I https://<public-bucket-domain>/packs/1/1.png     # 200 + immutable
curl -I https://<pending-bucket-domain>/pending/1.png    # 403/404
```

CORS, from the browser console on the live site:

```js
await fetch('https://<public>/packs/1/1.png', {mode:'cors'})   // must not throw
```

Then walk the full flow on a phone on mobile data: browse → copy → paste → favorite → go
offline → reload favorites.

## Files

- `docker-compose.yml`, deployment config
- `config/config.yaml`
- `README.md` (deployment notes)

## Hints

- Verify CORS with `fetch`, never with an `<img>` tag. This trips people every time.
- Run through spec §9's 16 success criteria one by one against production. That list is the
  definition of phase 1 being done.
