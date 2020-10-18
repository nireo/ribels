package commands

import "github.com/bwmarrin/discordgo"

// Github command just gives the link to the source code on github
func GithubCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	session.ChannelMessageSend(msg.ChannelID, "https://github.com/nireo/ribels")
}
