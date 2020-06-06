package db

import (
	"fmt"
	"log"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Bootstrap gorm mysql dialect
	"github.com/spf13/viper"
)

type config struct {
	host    string
	user    string
	pass    string
	dbName  string
	verbose bool
}

// NewDB configures and returns a DB connection
func NewDB() *gorm.DB {
	c := newConfig()
	db, err := gorm.Open("mysql", prepareConnectionString(c))
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
	)
}

func prepareConnectionString(c *config) string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", c.user, c.pass, c.host, c.dbName)
}

func newConfig() *config {
	return &config{
		host:    viper.GetString("db.host"),
		user:    viper.GetString("db.user"),
		pass:    viper.GetString("db.pass"),
		dbName:  viper.GetString("db.database"),
		verbose: viper.GetBool("db.verbose"),
	}
}
