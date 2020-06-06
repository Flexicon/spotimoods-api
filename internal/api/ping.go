package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type pingController struct{}

func newPing() Controller {
	return &pingController{}
}

func (c *pingController) Routes(g *echo.Group) {
	g.GET("/ping", func(c echo.Context) error {
		type response struct {
			Pong int64  `json:"pong"`
			Msg  string `json:"msg,omitempty"`
		}

		return c.JSON(http.StatusOK, response{
			Pong: time.Now().Unix(),
			Msg:  c.QueryParam("msg"),
		})
	})
}
