// Cross-cutting htmx behavior, bound once on document.

/**
 * Swap error responses (>= 400) into <body> instead of dropping them, so the
 * server-rendered error page is shown rather than a silent no-op.
 */
export function initErrorSwap(): void {
  document.body.addEventListener("htmx:beforeSwap", (evt) => {
    const detail = (evt as CustomEvent).detail as {
      xhr: XMLHttpRequest;
      shouldSwap: boolean;
      target: Element | null;
    };
    if (detail.xhr.status >= 400) {
      detail.shouldSwap = true;
      detail.target = htmx.find("body");
    }
  });
}

/**
 * Attach the CSRF token to every non-GET request. The token is rendered by the
 * CSRFToken component into <div id="csrf-token" data-token> inside <body> (so
 * hx-boost swaps keep it current) and read here at request time — no token is
 * ever interpolated into a <script>.
 */
export function initCsrf(): void {
  document.body.addEventListener("htmx:configRequest", (evt) => {
    const detail = (evt as CustomEvent).detail as {
      verb: string;
      parameters: Record<string, string>;
    };
    if (detail.verb === "get") return;
    const token = document.getElementById("csrf-token")?.dataset.token;
    if (token) detail.parameters["csrf"] = token;
  });
}
