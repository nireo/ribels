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

type LeagueUser struct {
	gorm.Model
	Username  string `json:"league_user"`
	Server    string `json:"server"`
	DiscordID string `json:"discord_id"`
	Region    string `json:"region"`
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

	if err != nil {
		log.Fatal(err)
	}

	// migrate models
	if err := database.AutoMigrate(&User{}, &LeagueUser{}); err != nil {
		log.Fatal(err)
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
