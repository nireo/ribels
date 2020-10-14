package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordID string `json:"discord_id"`
	OsuName   string `json:"osu_name"`
}

var db *gorm.DB

func InitDatabase() {
	user := os.Getenv("db_username")
	port := os.Getenv("db_port")
	host := os.Getenv("db_host")
	dbName := os.Getenv("db_name")

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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Problem loading environment file")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("Cannot create a bot instance")
	}

	dg.AddHandler(messageHandler)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection")
	}

	log.Print("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// tokenize the input
	args := strings.Split(msg.Content, " ")
	if args[0] == "$set" {
		if len(args) == 1 {
			session.ChannelMessageSend(msg.ChannelID, "No username provided")
			return
		}
	}
}
