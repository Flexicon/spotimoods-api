package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// ViperInit loads a viper config file and sets up needed defaults
func ViperInit() {
	// Viper init
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	// Defaults
	viper.SetDefault("port", 80)

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
