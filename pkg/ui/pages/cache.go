package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/tung-dnt/pagoda/pkg/ui"
	"github.com/tung-dnt/pagoda/pkg/ui/forms"
	"github.com/tung-dnt/pagoda/pkg/ui/layouts"
)

func UpdateCache(ctx echo.Context, form *forms.Cache) error {
	r := ui.NewRequest(ctx)
	r.Title = "Set a cache entry"

	return r.Render(layouts.Primary, form.Render(r))
}
