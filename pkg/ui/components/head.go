package components

import (
	"strings"

	"github.com/tung-dnt/pagoda/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func JS() Node {
	return Group{
		// An ES module so esbuild's split chunks load via import(); modules are deferred by default.
		Script(Type("module"), Src(ui.StaticFile("js/app.js"))),
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
