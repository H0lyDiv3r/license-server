package db

import (
	"fmt"
	"license-server/internals/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(&domain.User{}, &domain.License{})
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
