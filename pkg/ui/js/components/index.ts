// Registry for reusable Alpine components. Each component lives in its own file
// exporting a register function; wire it up here with an import plus a call
// inside the single alpine:init listener. Reference a registered component by
// name in gomponents markup, e.g. Div(Attr("x-data", "avatarEditor"), ...).
// Alpine's MutationObserver initializes them automatically, including on
// hx-boost swaps, so no manual re-init is needed after navigation.
//
// Split by concern: components/ = Alpine UI state; ../lib = framework-free
// data/platform logic. For heavy imperative widgets (charts, canvas,
// third-party libs), give the component init()/destroy() so setup and teardown
// are lifecycle-managed across swaps.
//
// Code splitting: components register eagerly (they are tiny), but any npm
// package or ../lib module only some pages need (IndexedDB, upload helpers,
// charts, croppers) must be loaded with import() from the component that uses
// it, so esbuild emits it as a separate js/chunks/ file. Never import such a
// package statically anywhere, or it lands back in app.js. Start the load in
// init() (and clean up in destroy()); for user-gesture APIs (clipboard, share,
// fullscreen) it must be loaded before the tap. See PAGODA.md "Page-level packages".

export function registerComponents(): void {
  document.addEventListener("alpine:init", () => {
    // Register components here — one import above, one call per component.
  });
}
