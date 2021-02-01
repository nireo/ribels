package commands

import "github.com/bwmarrin/discordgo"

// GithubCommandHandler just gives the link to the source code of the bot on github.
func GithubCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	_, _ = session.ChannelMessageSend(msg.ChannelID, "https://github.com/nireo/ribels")
}
