package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func BalanceCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// check if the user has a balance already
	db := utils.GetDatabase()
	var economyUser utils.EconomyUser
	if err := db.Where(&utils.EconomyUser{DiscordID: msg.Author.ID}).First(&economyUser).Error; err == nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Your balance is `%d`", economyUser.Balance))
		return
	}

	// if the user doesn't have a balance, start the user out with 100 coins
	newUser := &utils.EconomyUser{
		Balance: 100,
		DiscordID: msg.Author.ID,
	}

	db.Create(&newUser)
	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("Created a bank account for you!\nCurrent balance is `100`!"))
}
