package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/backlite/ui"
	"github.com/tung-dnt/pagoda/pkg/middleware"
	"github.com/tung-dnt/pagoda/pkg/routenames"
	"github.com/tung-dnt/pagoda/pkg/services"
)

type Admin struct {
	backlite *ui.Handler
}

func init() {
	Register(new(Admin))
}

func (h *Admin) Init(c *services.Container) error {
	var err error
	h.backlite, err = ui.NewHandler(ui.Config{
		DB:           c.TasksDatabase,
		BasePath:     "/admin/tasks",
		ItemsPerPage: 25,
		ReleaseAfter: c.Config.Tasks.ReleaseAfter,
	})
	return err
}

func (h *Admin) Routes(g *echo.Group) {
	ag := g.Group("/admin", middleware.RequireAdmin)

	tasks := ag.Group("/tasks")
	tasks.GET("", h.Backlite(h.backlite.Running)).Name = routenames.AdminTasks
	tasks.GET("/succeeded", h.Backlite(h.backlite.Succeeded))
	tasks.GET("/failed", h.Backlite(h.backlite.Failed))
	tasks.GET("/upcoming", h.Backlite(h.backlite.Upcoming))
	tasks.GET("/task/:id", h.Backlite(h.backlite.Task))
	tasks.GET("/completed/:id", h.Backlite(h.backlite.TaskCompleted))
}

func (h *Admin) Backlite(handler func(http.ResponseWriter, *http.Request) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if id := c.Param("id"); id != "" {
			c.Request().SetPathValue("task", id)
		}
		return handler(c.Response().Writer, c.Request())
	}
}
