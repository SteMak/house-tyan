package database

import (
	"log"

	"github.com/SteMak/ani-helper/workerTools/config"
	"github.com/jinzhu/gorm"
	// driver db
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// Records is an object of records in database
	Records *records

	database *gorm.DB
)

// Init intialies database
func Init() {
	var err error
	database, err = gorm.Open("postgres", config.PostgresURI)
	if err != nil {
		log.Fatal(err)
	}
	database.DB().Ping()

	database.AutoMigrate(&Record{})
	Records = &records{
		db: database,
	}
}

// Close closes database connection
func Close() {
	database.Close()
}
