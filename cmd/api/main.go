package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/api"
	"github.com/flexicon/spotimoods-go/internal/config"
	"github.com/flexicon/spotimoods-go/internal/db"
	"github.com/flexicon/spotimoods-go/internal/spotify"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()

	d := db.NewDB()
	h := &http.Client{Timeout: 5 * time.Second}

	repos := db.NewRepositoryProvider(d)
	spot := spotify.NewClient(h)

	services := internal.NewServiceProvider(repos, spot)

	api.InitRoutes(e, api.Options{
		Services: services,
	})

	port := viper.GetInt("web.port")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func init() {
	config.ViperInit()
}
