package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"math/rand"
	"strconv"
	"time"
)

func CoinflipCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) == 1 {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "You need to provide a wager!")
		return
	}

	wager, err := strconv.Atoi(args[1])
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Cannot parse wager as a number!")
		return
	}

	// check if the user has a balance already
	db := utils.GetDatabase()
	var economyUser utils.EconomyUser
	if err := db.Where(&utils.EconomyUser{DiscordID: msg.Author.ID}).First(&economyUser).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"No bank account found, create a bank account with `;balance`!")
		return
	}

	if economyUser.Balance < int64(wager) {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "You don't have sufficient funds!")
		return
	}

	if wager < 1 {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "You have to bet something!")
		return
	}

	var won bool
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(200)
	if randomNumber <= 200 && randomNumber >= 100 {
		economyUser.Balance += int64(wager)
		won = true
	} else {
		economyUser.Balance -= int64(wager)
		won = false
	}

	if won {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("You won!\nCurrent balance is: `%d`.", economyUser.Balance))
	} else {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("You lost!\nCurrent balance is: `%d`", economyUser.Balance))
	}

	db.Save(&economyUser)
}