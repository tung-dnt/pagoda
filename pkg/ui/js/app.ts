// Entry point for all UI JavaScript. Bundled by esbuild into
// public/static/js/app.js (see `make js`) and loaded once in <head> as an ES module.
//
// htmx and Alpine are bundled in here from npm (see package.json) rather than
// loaded from a CDN, so package.json is the single source of truth for their
// versions and the app ships fully self-contained (no third-party runtime deps).
//
// Everything here binds to `document` or registers with htmx/Alpine, so it stays
// active across hx-boost navigations without re-running per page. To add page-level
// interactivity, prefer an Alpine component in ./components — not a per-page script.

import Alpine from "alpinejs";
import htmx from "htmx.org";

import { registerComponents } from "./components";
import { initCsrf, initErrorSwap } from "./lib/htmx";

// Expose the globals the CDN builds used to create, for parity (inline hx-on
// handlers, Alpine plugins, browser devtools) and for our lib modules that
// reference the `htmx` global.
window.htmx = htmx;
window.Alpine = Alpine;

// Register Alpine components before Alpine.start() dispatches `alpine:init`.
registerComponents();

function boot(): void {
  initErrorSwap();
  initCsrf();
}

// Module scripts are deferred, so parsing is usually done — but guard for safety.
if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", boot);
} else {
  boot();
}

// Alpine's npm build does not auto-start (unlike the CDN build) — start it
// explicitly. It scans the already-parsed DOM and its observer handles later
// hx-boost swaps.
Alpine.start();
