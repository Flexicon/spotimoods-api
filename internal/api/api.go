package api

import (
	"fmt"
	"net/http"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Controller defines all methods for the API controllers
type Controller interface {
	Routes(g *echo.Group)
}

// Options for API routes
type Options struct {
	Services *internal.ServiceProvider
}

// ErrResponse for generic API error messages
type ErrResponse struct {
	Msg string `json:"message"`
}

// InitRoutes setup router, middleware and mounts all controllers
func InitRoutes(e *echo.Echo, opts Options) {
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	base := e.Group("")
	newHome().Routes(base)
	newLogin(opts.Services).Routes(base)

	api := e.Group("/api")
	newPing().Routes(api)
	newUser(opts.Services).Routes(api)
	newMood(opts.Services).Routes(api)
}

func notFound(c echo.Context, resource string) error {
	msg := "not found"
	if resource != "" {
		msg = fmt.Sprintf("%s %s", resource, msg)
	}

	return c.JSON(http.StatusNotFound, ErrResponse{Msg: msg})
}
