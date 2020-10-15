package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nireo/ribels/handlers"
	"github.com/nireo/ribels/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Problem loading environment file")
	}

	utils.InitDatabase()
	utils.InitApiKey()

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("Cannot create a bot instance")
	}

	dg.AddHandler(handlers.MessageHandler)
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
