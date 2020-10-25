package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func SetLeagueCommandHandler(session *discordgo.Session,
	msg *discordgo.MessageCreate, args []string) {
	// check that an actual username is provided
	if len(args) < 3 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Not enough arguments provided! Provide them: USERNAME, SERVER")
		return
	}

	db := utils.GetDatabase()

	var leagueUser utils.LeagueUser
	if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&leagueUser).Error; err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"A league account has already been linked with this account!")
		return
	}

	leagueUsername := utils.FormatName(args[1:])

	// create the new user model
	newUser := &utils.LeagueUser{
		Username:  leagueUsername,
		DiscordID: msg.Author.ID,
		Server:    args[len(args)-1],
	}

	db.Create(&newUser)
	_, _ = session.ChannelMessageSend(msg.ChannelID, "The provided is user is now linked to the discord!")
}
