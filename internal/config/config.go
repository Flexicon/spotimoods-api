package config

import (
	"flag"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// ViperInit loads a viper config file and sets up needed defaults
func ViperInit() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	// Defaults
	viper.SetDefault("port", 80)

	initFlags()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Failed to load config file: ", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalln("Viper error: ", err)
		}
	}
}

func initFlags() {
	workerPtr := flag.Bool("worker", false, "whether to run as a background worker")
	flag.Parse()

	viper.Set("worker", *workerPtr)
}
