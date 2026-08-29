package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/tung-dnt/meme-app/pkg/pager"
	"github.com/tung-dnt/meme-app/pkg/routenames"
	"github.com/tung-dnt/meme-app/pkg/services"
	"github.com/tung-dnt/meme-app/pkg/ui/models"
	"github.com/tung-dnt/meme-app/pkg/ui/pages"
)

type Pages struct{}

func init() {
	Register(new(Pages))
}

func (h *Pages) Init(c *services.Container) error {
	return nil
}

// Routes registers the public application pages on the public group so they are edge-cacheable.
// Neither page contains a CSRF-protected form; a logged-in visitor still renders personalized
// (Session/LoadAuthenticatedUser run in the shared base) and bypasses the cache at the edge.
func (h *Pages) Routes(_, pub *echo.Group) {
	pub.GET("/", h.Home).Name = routenames.Home
	pub.GET("/about", h.About).Name = routenames.About
}

func (h *Pages) Home(ctx echo.Context) error {
	pgr := pager.NewPager(ctx, 4)

	return pages.Home(ctx, &models.Posts{
		Posts: h.fetchPosts(&pgr),
		Pager: pgr,
	})
}

// fetchPosts is a mock example of fetching posts to illustrate how paging works.
func (h *Pages) fetchPosts(pager *pager.Pager) []models.Post {
	pager.SetItems(20)
	posts := make([]models.Post, 20)

	for k := range posts {
		posts[k] = models.Post{
			ID:    k + 1,
			Title: fmt.Sprintf("Post example #%d", k+1),
			Body:  fmt.Sprintf("Lorem ipsum example #%d ddolor sit amet, consectetur adipiscing elit. Nam elementum vulputate tristique.", k+1),
		}
	}
	return posts[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}

func (h *Pages) About(ctx echo.Context) error {
	return pages.About(ctx)
}
