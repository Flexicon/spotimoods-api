package db

import (
	"log"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Bootstrap gorm mysql dialect
	"github.com/spf13/viper"
)

type config struct {
	connectionURI string
	verbose       bool
}

// NewDB configures and returns a DB connection
func NewDB() *gorm.DB {
	c := newConfig()
	db, err := gorm.Open("mysql", c.connectionURI)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	if c.verbose {
		db.LogMode(true)
	}

	autoMigrate(db)

	return db
}

func autoMigrate(d *gorm.DB) {
	d.AutoMigrate(
		&internal.User{},
		&internal.SpotifyToken{},
		&internal.Mood{},
	)
}

func newConfig() *config {
	return &config{
		connectionURI: viper.GetString("db.uri"),
		verbose:       viper.GetBool("db.verbose"),
	}
}
