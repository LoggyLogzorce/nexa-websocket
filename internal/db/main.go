package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var connect = "host=localhost port=5432 user=postgres password=1234 dbname=nexa_warp_db sslmode=disable"

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	log.Println("Connected to the database")
	return db
}
