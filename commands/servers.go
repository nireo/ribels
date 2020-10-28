package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

// print all the available league of legends servers with their shorthand name
func ServersCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	message := ""
	for server := range utils.ValidRegions {
		message += server+"\n"
	}

	_, _ = session.ChannelMessageSend(msg.ChannelID, message)
}
