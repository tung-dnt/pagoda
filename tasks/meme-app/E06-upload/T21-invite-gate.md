# T21 · Upload form + invite gate + rate limit

| | |
|---|---|
| **Epic** | [E06 · Upload](./README.md) |
| **Size** | M — ~2.5 h |
| **Depends on** | T05 |
| **Unlocks** | T22 |
| **Spec** | §6.5 |
| **Status** | ☐ not started |

## Problem

Build `GET /upload`: the form, the invite-code check, and the rate limiter. No file
handling yet — this task is purely about who is allowed to reach the upload endpoint.

Splitting the gate from the ingestion keeps the security-sensitive part reviewable on its
own.

## Given

- `pkg/ui/forms/file.go` — the form component pattern
- `pkg/ui/forms/login.go` — a form with validation
- `pkg/form` — form binding and validation helpers
- Echo's `middleware.RateLimiter`

## Constraints

- Invite code lives in config (`app.inviteCode`), overridden by env. **Never committed.**
- Compare with `subtle.ConstantTimeCompare`. A plain `==` on a secret leaks its length and
  prefix through timing.
- Rate limit: **10 attempts per hour per IP**, applied to the failure path. Without this
  the code is brute-forceable at leisure.
- Behind a proxy, `ctx.RealIP()` trusts `X-Forwarded-For`. Make sure the deployment sets it
  correctly, or the limiter keys every request to the same address and locks everyone out.
- The pack-name field is shown **only for zip uploads** — a single image goes to the global
  "Singles" pack and has no name of its own.
- CSRF is already applied globally in `router.go`; include the token via the existing
  `CSRF(r)` helper.
- A wrong code gives a generic failure. Do not reveal whether the code was close.

## Acceptance

- [ ] `GET /upload` renders: invite code, optional pack name, tags, file input
- [ ] A wrong code is rejected with a generic message
- [ ] The correct code proceeds to the (stubbed) upload path
- [ ] The 11th failed attempt within an hour is rate-limited
- [ ] The code is compared in constant time
- [ ] No code value appears in logs, ever
- [ ] The form carries a CSRF token

## Verify

```
make test && make run
```

```
for i in $(seq 1 12); do curl -s -o /dev/null -w "%{http_code}\n" \
  -X POST localhost:8000/upload -F invite=wrong; done
```

Expect 429 by the 11th. Then `grep -ri "invite" logs/` — expect no secret values.

## Files

- `pkg/handlers/upload.go`, `pkg/handlers/upload_test.go`
- `pkg/ui/pages/upload.go`, `pkg/ui/forms/upload.go`
- `config/config.go`, `config/config.yaml`
- `pkg/routenames/names.go`

## Hints

- `subtle.ConstantTimeCompare` returns 1 on equal. Compare fixed-length hashes of both
  sides if you want to avoid leaking length as well.
- Write the rate-limit test with a fake clock; a test that actually waits an hour is not a
  test you will run.
