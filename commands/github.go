package commands

import "github.com/bwmarrin/discordgo"

func GithubCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	session.ChannelMessageSend(msg.ChannelID, "https://github.com/nireo/ribels")
}
