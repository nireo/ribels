package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func ReadyHandler(session *discordgo.Session, event *discordgo.Ready) {
	if err := session.UpdateStatus(0, "$help"); err != nil {
		log.Fatalln("Error settings discord status")
	}
}
