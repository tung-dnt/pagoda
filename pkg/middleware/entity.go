package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/tung-dnt/pagoda/pkg/context"
	pgdb "github.com/tung-dnt/pagoda/pkg/postgres/db"

	"github.com/labstack/echo/v4"
)

// LoadUser loads the user based on the ID provided as a path parameter.
func LoadUser(db *pgdb.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := strconv.ParseInt(c.Param("user"), 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			u, err := db.GetUser(c.Request().Context(), userID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return echo.NewHTTPError(http.StatusNotFound)
				}

				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error querying user: %v", err),
				)
			}

			c.Set(context.UserKey, &u)

			return next(c)
		}
	}
}
