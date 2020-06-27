package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/api"
	"github.com/flexicon/spotimoods-go/internal/config"
	"github.com/flexicon/spotimoods-go/internal/db"
	"github.com/flexicon/spotimoods-go/internal/queue"
	"github.com/flexicon/spotimoods-go/internal/spotify"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()

	d := db.NewDB()
	h := &http.Client{Timeout: 5 * time.Second}

	// Init all app services
	repos := db.NewRepositoryProvider(d)
	spot := spotify.NewClient(h, repos)
	qs, err := queue.Setup()
	if err != nil {
		log.Fatalln(err)
	}

	// Setup main service provider
	services := internal.NewServiceProvider(repos, spot, qs)

	// Setup API and queue consumers
	api.InitRoutes(e, api.Options{
		Services: services,
	})
	go func() {
		qh := queue.NewHandler(services)
		log.Fatalln(queue.Listen(qs, qh))
	}()

	// Test queue connection with a ping message
	go func() {
		if err := qs.Ping("ping"); err != nil {
			log.Fatal(err)
		}
	}()

	// Start up web server
	port := viper.GetInt("web.port")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func init() {
	config.ViperInit()
}
