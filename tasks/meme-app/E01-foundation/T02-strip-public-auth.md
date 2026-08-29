# T02 · Strip public registration

| | |
|---|---|
| **Epic** | [E01 · Foundation](./README.md) |
| **Size** | S — ~1 h |
| **Depends on** | None |
| **Unlocks** | T03 |
| **Spec** | §6.6 |
| **Status** | ☐ not started |

## Problem

The boilerplate ships a full public auth system. This app has exactly one account, created
by `make admin`. Every public auth route is attack surface for a feature that does not
exist, so delete it rather than hide it.

Keep: `/login`, `/logout`, and the session machinery. Remove everything else.

## Given

- `pkg/handlers/auth.go` — handlers for register, forgot password, reset password, verify email
- `pkg/ui/pages/auth.go` — the corresponding pages
- `pkg/ui/forms/register.go`, `forgot_password.go`, `reset_password.go`
- `pkg/routenames/names.go` — the route name constants
- `pkg/ui/components/nav.go` — nav links pointing at them
- `pkg/services/auth.go` — `AuthClient`; keep it, it backs login

## Constraints

- Delete route *registrations*, not just the nav links. An unlinked route is still routable.
- `password_tokens` and the email-verification code paths go with it, but **leave the table
  in place** — dropping it is a migration this task does not need and T04 will not touch it.
- `make admin` must still work afterwards. It calls `c.Auth.HashPassword` and
  `c.Queries.CreateUser` — neither should be affected.
- Existing tests reference the removed routes; delete those test cases rather than skipping.

## Acceptance

- [ ] `GET /register`, `/forgot-password`, `/reset-password` all return 404
- [ ] `GET /login` renders and a valid admin login still lands authenticated
- [ ] No dead route-name constants remain in `pkg/routenames/names.go`
- [ ] Nav shows no registration link to an anonymous visitor
- [ ] `make admin email=test@local` still creates a working admin

## Verify

```
make test
make build
make admin email=spec-check@local     # then log in with the printed password
```

## Files

- `pkg/handlers/auth.go`, `pkg/handlers/auth_test.go`
- `pkg/ui/pages/auth.go`
- `pkg/ui/forms/register.go`, `forgot_password.go`, `reset_password.go` (delete)
- `pkg/routenames/names.go`
- `pkg/ui/components/nav.go`

## Hints

- Grep for each constant name before deleting it — `grep -rn RegisterSubmit pkg/`.
- The email templates in `pkg/ui/emails/auth.go` become unreachable too.
