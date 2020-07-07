package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type homeController struct{}

func newHome() Controller {
	return &homeController{}
}

func (c *homeController) Routes(g *echo.Group) {
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Spotimoods™")
	})
}
