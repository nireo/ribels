package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"strings"
)

func SetLeagueCommandHandler(session *discordgo.Session,
	msg *discordgo.MessageCreate, args []string) {
	// check that an actual username is provided
	if len(args) != 3 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Invalid arguments! Provide them: USERNAME, SERVER")
		return
	}

	db := utils.GetDatabase()

	var leagueUser utils.LeagueUser
	if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&leagueUser).Error; err == nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"A league account has already been linked with this account!")
		return
	}

	if _, ok := utils.ValidRegions[strings.ToLower(args[2])]; !ok {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Problem parsin league server")
		return
	}

	// create the new user model
	newUser := &utils.LeagueUser{
		Username:  args[1],
		DiscordID: msg.Author.ID,
		Region:    args[len(args)-1],
	}

	db.Create(&newUser)
	_, _ = session.ChannelMessageSend(msg.ChannelID, "The provided is user is now linked to the discord!")
}
