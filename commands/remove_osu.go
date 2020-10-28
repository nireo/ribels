package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func RemoveOsuCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	var user utils.User
	db := utils.GetDatabase()

	if err := db.Where(&utils.User{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"An osu! account is not linked to your discord profile")
		return
	}

	db.Delete(&user)
	_, _ = session.ChannelMessageSend(msg.ChannelID,
		"An osu! account has been unlinked from your discord id!")
}
