package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func ResetBalance(session *discordgo.Session, msg *discordgo.MessageCreate) {
	db := utils.GetDatabase()
	var economyUser utils.EconomyUser
	if err := db.Where(&utils.EconomyUser{DiscordID: msg.Author.ID}).First(&economyUser).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "No back account found! Create with `balance`")
		return
	}

	economyUser.Balance = 100
	db.Save(&economyUser)
}
