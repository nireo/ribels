package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

// This command makes a link between a discord id and a osu username, so we can use commands without arguments
func SetCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	// Check that the name isn't empty
	if len(args) == 1 {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "No username provided!")
		return
	}

	// format the name into a better format: "install gentoo" => "install_gentoo"
	osuName := utils.FormatName(args[1:])

	db := utils.GetDatabase()
	var user utils.User

	// err == nil, because if there was no problem finding an user, that means that the user already exists
	if err := db.Where(&utils.User{OsuName: osuName}).First(&user).Error; err == nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "User already in database")
		return
	}

	// create the new model to the database
	newUser := &utils.User{
		OsuName:   osuName,
		DiscordID: msg.Author.ID,
	}

	// save the model
	db.Create(&newUser)
	_, _ = session.ChannelMessageSend(msg.ChannelID, "Username saved into database")
}
