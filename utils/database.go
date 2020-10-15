package utils

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordID string `json:"discord_id"`
	OsuName   string `json:"osu_name"`
}

var db *gorm.DB

func GetDatabase() *gorm.DB {
	return db
}

func InitDatabase() {
	user := os.Getenv("DATABASE_USER")
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	dbName := os.Getenv("DATABASE_NAME")

	// Load database
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbName),
	}), &gorm.Config{})

	// migrate models
	database.AutoMigrate(&User{})

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	db = database
}

func CheckIfSet(userID string) (User, error) {
	var user User
	if err := db.Where(&User{DiscordID: userID}).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}