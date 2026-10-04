package components

import (
	"github.com/tung-dnt/pagoda/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// CSRFToken exposes the CSRF token to js/lib/htmx.ts, which attaches it to every non-GET htmx
// request. It lives in <body> rather than a <head> meta tag on purpose: hx-boost swaps only the
// body, so a head tag would go stale when navigating from a token-less public page (see the public
// route group) to a form page. Rendered only when the CSRF middleware populated a token, so cached
// public pages never carry one.
func CSRFToken(r *ui.Request) Node {
	return If(len(r.CSRF) > 0, Div(ID("csrf-token"), Class("hidden"), Data("token", r.CSRF)))
}

// HxBoost enables htmx boosting so anchors/forms navigate via AJAX body swaps
// instead of full page loads. The behavior it relies on (CSRF, error swapping)
// lives in pkg/ui/js/lib/htmx.ts.
func HxBoost() Node {
	return Attr("hx-boost", "true")
}
