package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nireo/ribels/handlers"
	"github.com/nireo/ribels/utils"
)

func main() {
	// Load all the environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Problem loading environment file")
	}

	// init the database
	utils.InitDatabase()

	// init the osu api
	utils.InitApiKey()

	// init the league of legends api
	utils.InitClient()

	// check if logging is enabled, and set the logging in the message handler
	status, _ := strconv.ParseBool(os.Getenv("LOGGING"))
	handlers.SetLogging(status)

	// Create discord instance
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("Cannot create a bot instance")
	}

	// add the message handler, which handles all the commands
	dg.AddHandler(handlers.MessageHandler)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	if err := dg.Open(); err != nil {
		log.Fatal("Error opening connection")
	}

	log.Print("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := dg.Close(); err != nil {
		log.Fatal(err)
	}
}
