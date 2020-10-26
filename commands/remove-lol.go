package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func RemoveLolCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	db := utils.GetDatabase()

	var user utils.LeagueUser
	if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "No link between league user exists")
		return
	}

	db.Delete(&user)
}