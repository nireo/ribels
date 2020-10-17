package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func SetCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) == 1 {
		session.ChannelMessageSend(msg.ChannelID, "No username provided!")
		return
	}

	osuName := utils.FormatName(args[1:])
	db := utils.GetDatabase()

	var user utils.User
	if err := db.Where(&utils.User{OsuName: osuName}).First(&user).Error; err != nil {
		session.ChannelMessageSend(msg.ChannelID, "User already in database")
		return
	}

	newUser := &utils.User{
		OsuName:   osuName,
		DiscordID: msg.Author.ID,
	}

	db.Create(&newUser)
	session.ChannelMessageSend(msg.ChannelID, "Username saved into database")
}
