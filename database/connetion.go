package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	Conn *gorm.DB
)

func Connect() {
	dsn := fmt.Sprintf("host=45.135.56.198 user=admin password=%s dbname=edward port=5432", os.Getenv("DB_PASS"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	Conn = db

	err = db.AutoMigrate(&User{}, &Session{}, &Guild{}, &Permissions{})
	if err != nil {
		log.Println(err.Error())
	}
}
