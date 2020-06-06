package api

import (
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Controller defines all methods for the API controllers
type Controller interface {
	Routes(g *echo.Group)
}

// Options for InitRoutes.
type Options struct {
	Services *internal.ServiceProvider
}

// InitRoutes setup router, middleware and mounts all controllers
func InitRoutes(e *echo.Echo, opts Options) {
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{},
		AllowMethods: []string{echo.GET},
	}))

	base := e.Group("")
	newLogin(opts.Services).Routes(base)

	api := e.Group("/api")
	newPing().Routes(api)
	newAuth(opts.Services).Routes(api)
}
