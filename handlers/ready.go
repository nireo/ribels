package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

// ReadyHandler handles setting the status displayed under the bot name on the sidebar of the server.
func ReadyHandler(session *discordgo.Session, event *discordgo.Ready) {
	if err := session.UpdateStatus(0, utils.FormatCommand("help")); err != nil {
		log.Fatalln("Error settings discord status")
	}
}
