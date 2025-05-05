package config

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BuildDB() *gorm.DB {
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")
	sqlCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DbHost,
		DbPort,
		DbUser,
		DbPassword,
		DbName,
	)

	db, err := gorm.Open(postgres.Open(sqlCfg), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to database postgres")
		log.Fatal("connection error:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	return db
}
