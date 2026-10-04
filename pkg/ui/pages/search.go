package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/tung-dnt/pagoda/pkg/ui"
	"github.com/tung-dnt/pagoda/pkg/ui/layouts"
	"github.com/tung-dnt/pagoda/pkg/ui/models"
	. "maragu.dev/gomponents"
)

func SearchResults(ctx echo.Context, results []*models.SearchResult) error {
	r := ui.NewRequest(ctx)

	g := make(Group, len(results))
	for i, result := range results {
		g[i] = result.Render()
	}

	return r.Render(layouts.Primary, g)
}
