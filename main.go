package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type OsuUserResponse struct {
	Username      string `json:"username"`
	Playcount     string `json:"playcount"`
	RankedScore   string `json:"ranked_score"`
	PPRank        string `json:"pp_rank"`
	Level         string `json:"level"`
	Accuracy      string `json:"accuracy"`
	Country       string `json:"country"`
	SecondsPlayed string `json:"total_seconds_played"`
}

type OsuAPI struct {
	Key string
}

func (osuApi *OsuAPI) GetUser(username string) (OsuUserResponse, error) {
	var osuUser OsuUserResponse
	response, err := http.Get("api_call")
	if err != nil {
		return osuUser, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return osuUser, err
	}

	if err := json.Unmarshal(body, &osuUser); err != nil {
		return osuUser, err
	}

	return osuUser, nil
}

var db *gorm.DB
var api *OsuAPI

func CheckIfSet(userID string) (User, error) {
	var user User
	if err := db.Where(&User{DiscordID: userID}).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Problem loading environment file")
	}

	InitDatabase()

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("Cannot create a bot instance")
	}

	api = &OsuAPI{Key: os.Getenv("OSU_KEY")}
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

		osu_name := strings.Join(args[1:], "_")
		// check if name already in database
		var user User
		if err := db.Where(&User{OsuName: osu_name}).First(&user).Error; err == nil {
			session.ChannelMessageSend(msg.ChannelID, "User already in database")
			return
		}

		// insert into database
		newUser := &User{
			DiscordID: msg.Author.ID,
			OsuName:   osu_name,
		}

		db.Create(&newUser)

		session.ChannelMessageSend(msg.ChannelID, "Saved user in database")
	}

	if args[0] == "$osu" {
		user, err := CheckIfSet(msg.Author.ID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Not set in database")
			return
		}

		session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is set", user.OsuName))
	}

	if args[0] == "$top" {
		user, err := CheckIfSet(msg.Author.ID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Not set in database")
			return
		}

		session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is set", user.OsuName))
	}
}
