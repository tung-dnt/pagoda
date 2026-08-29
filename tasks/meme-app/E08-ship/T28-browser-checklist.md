# T28 · Manual browser verification

| | |
|---|---|
| **Epic** | [E08 · Ship](./README.md) |
| **Size** | S — ~1.5 h |
| **Depends on** | T26 |
| **Unlocks** | T29 |
| **Spec** | §7 |
| **Status** | ☐ not started |

## Problem

The client-side layer — clipboard, IndexedDB, favorites — has **no automated tests**. That
was a deliberate, time-driven cut and it is the single largest quality risk in the plan.
This task is its mitigation: a written checklist, actually executed, on real browsers.

Write it down so it is repeatable, not a one-off poke.

## Given

- Everything from T10–T15
- `docs/spike-clipboard.md` from T01 — the quirks you already found

## Constraints

- Four targets, and **iOS Safari is not optional** — it is where clipboard behaviour
  differs most and where the spike found its surprises:
  - Chrome desktop
  - Firefox desktop
  - Safari desktop
  - **iOS Safari on a real device**
- Test over **HTTPS**. Clipboard needs a secure context, and `localhost` behaving is not
  evidence.
- Include the paste targets that matter, not just the copy: Slack, Discord, iOS Messages,
  WhatsApp or Zalo.
- Include a private-window run — IndexedDB may be denied, and the app must degrade rather
  than break.
- Record actual results, not just checkboxes. A note reading "Firefox: animated copy needs
  two taps" is worth more than a tick.

## Acceptance

- [ ] `docs/browser-checklist.md` exists with every case enumerated
- [ ] The full checklist has been run on all four targets, results recorded
- [ ] Copy (both paths) verified against at least three paste targets
- [ ] Offline favorites verified on desktop and on iOS
- [ ] Private-window degradation verified
- [ ] Every failure found is either fixed or filed with a written decision to accept it

## Verify

The checklist itself is the verification. It is done when every row has a real result.

## Files

- `docs/browser-checklist.md`

## Hints

- Cover at minimum: copy static, copy animated, favorite, unfavorite, favorite a pack,
  offline reload of `/favorites`, blob eviction at the cap, search, tag filter, upload,
  approve.
- Anything you find here that is cheap to guard with a Go test, add the test — some
  client-visible bugs have server-side causes.
