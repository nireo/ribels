package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func VoteCommandHandler(session *discordgo.Session, msg *discordgo.Message, args []string) {
	question := strings.Join(args[1:], " ")
	if question == "" {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Cannot vote on a empty question!")
		return
	}

	var messageEmbed discordgo.MessageEmbed
	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
