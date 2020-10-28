package handlers

import "github.com/bwmarrin/discordgo"

func ReadyHandler(session *discordgo.Session, event *discordgo.Ready) {
	session.UpdateStatus(0, "$help")
}
