package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/api"
	"github.com/flexicon/spotimoods-go/internal/cache"
	"github.com/flexicon/spotimoods-go/internal/config"
	"github.com/flexicon/spotimoods-go/internal/db"
	"github.com/flexicon/spotimoods-go/internal/queue"
	"github.com/flexicon/spotimoods-go/internal/spotify"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	config.ViperInit()

	d := db.NewDB()
	h := &http.Client{Timeout: 5 * time.Second}

	// Init all app services
	repos := db.NewRepositoryProvider(d)
	qs, err := queue.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	cs, err := cache.NewCache()
	if err != nil {
		log.Fatalln(err)
	}
	spot := spotify.NewClient(h, repos, cs)

	// Setup main service provider
	services := internal.NewServiceProvider(repos, spot, qs, cs)

	// Queue consumers and test queue connection with a ping message
	go setupQueueListener(services)
	go pingQueue(qs)

	// Setup web server if not running as a background worker
	if !viper.GetBool("worker") {
		e := echo.New()
		api.InitRoutes(e, api.Options{
			Services: services,
		})

		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("port"))))
	}

	// Run worker until system interrupt signal is received
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func setupQueueListener(services *internal.ServiceProvider) {
	qh := queue.NewHandler(services)
	log.Fatalln(queue.Listen(services.Queue().(*queue.Service), qh))
}

func pingQueue(qs *queue.Service) {
	if err := qs.Ping("ping"); err != nil {
		log.Fatal(err)
	}
}
