package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"strconv"
)

func PayCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) <= 2 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Need more arguments, for example: `pay 100 @ribels#0001`")
		return
	}

	if len(msg.Mentions) == 0 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"You need to mention the user you want to pay!")
		return
	}

	payment, err := strconv.Atoi(args[1])
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Cannot parse payment as a number!")
		return
	}

	if payment < 1 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"You need to pay a positive amount!")
		return
	}

	if msg.Author.ID == msg.Mentions[0].ID {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"You cannot pay yourself :P")
		return
	}

	db := utils.GetDatabase()
	var user utils.EconomyUser
	if err := db.Where(&utils.EconomyUser{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Could not find a bank account, create one with `;balance`!")
		return
	}

	if user.Balance < int64(payment) {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"You don't have sufficient funds!")
		return
	}

	var toSend utils.EconomyUser
	if err := db.Where(&utils.EconomyUser{DiscordID: msg.Mentions[0].ID}).First(&toSend).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Could not find a bank account for `%s`", msg.Mentions[0].Username))
		return
	}

	toSend.Balance += int64(payment)
	user.Balance -= int64(payment)

	db.Save(&toSend)
	db.Save(&user)

	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("Sent `%d` to `%s`", payment, msg.Mentions[0].Username))
}
