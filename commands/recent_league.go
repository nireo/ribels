package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
)

func RecentLeagueCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var region string
	var user utils.LeagueUser

	db := utils.GetDatabase()
	if len(args) != 3 {
		if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting user from database!")
			return
		}

		region = user.Region
	} else {
		region = args[2]
	}

	validRegion, err := utils.CheckValidRegion(region)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Error parsing region")
		return
	}

	_ = utils.NewRiotClient(validRegion, os.Getenv("LEAGUE_API"))
}
