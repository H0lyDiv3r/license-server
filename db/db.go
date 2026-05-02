package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("migration failed: %s", err.Error())
	}

}

func Connect() *gorm.DB {

	usr := os.Getenv("POSTGRES_USER")
	passwd := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%s sslmode=disable", usr, passwd, db, port)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect to db", err.Error())
	}

	Migrate(gormDB)

	return gormDB

}
