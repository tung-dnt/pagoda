package components

import (
	"strings"

	"github.com/tung-dnt/meme-app/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func JS() Node {
	return Group{
		Script(Src("https://unpkg.com/htmx.org@2.0.0/dist/htmx.min.js"), Defer()),
		// Pinned to an exact version, not a `3.x.x` range: the service worker caches this
		// URL, and a floating range would serve changing content behind a stable cache key.
		Script(Src("https://unpkg.com/alpinejs@3.16.3/dist/cdn.min.js"), Defer()),
	}
}

func CSS() Node {
	return Link(
		Href(ui.StaticFile("main.css")),
		Rel("stylesheet"),
		Type("text/css"),
	)
}

func Metatags(r *ui.Request) Node {
	return Group{
		Meta(Charset("utf-8")),
		Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		Link(Rel("icon"), Href(ui.StaticFile("favicon.png"))),
		TitleEl(Text(r.Config.App.Name), If(r.Title != "", Text(" | "+r.Title))),
		If(r.Metatags.Description != "", Meta(Name("description"), Content(r.Metatags.Description))),
		If(len(r.Metatags.Keywords) > 0, Meta(Name("keywords"), Content(strings.Join(r.Metatags.Keywords, ", ")))),
	}
}
