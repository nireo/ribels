package utils

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User database model
type User struct {
	gorm.Model
	DiscordID string `json:"discord_id"`
	OsuName   string `json:"osu_name"`
}

// LeagueUser database model
type LeagueUser struct {
	gorm.Model
	Username  string `json:"league_user"`
	DiscordID string `json:"discord_id"`
	Region    string `json:"region"`
}

// EconomyUser data model
type EconomyUser struct {
	gorm.Model
	DiscordID string `json:"discord_id"`
	Balance   int64  `json:"balance"`
}

var db *gorm.DB

// GetDatabase returns a pointer to the local db variable.
func GetDatabase() *gorm.DB {
	return db
}

// InitDatabase sets up the database given a few parameters from the environment variables.
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
	if err := database.AutoMigrate(&User{}, &LeagueUser{}, &EconomyUser{}); err != nil {
		log.Fatal(err)
	}

	db = database
}

// CheckIfSet checks if a user already exists for a given discord_id
func CheckIfSet(userID string) (User, error) {
	var user User
	if err := db.Where(&User{DiscordID: userID}).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
