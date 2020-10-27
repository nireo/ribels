package commands

import (
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
	"strings"
)

func RecentLeagueCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var leagueName string
	var region string
	var user utils.LeagueUser

	db := utils.GetDatabase()
	if len(args) != 3 {
		if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting user from database!")
			return
		}

		leagueName = user.Username
		region = user.Region
	} else {
		leagueName = args[1]
		region = args[2]
	}

	var reg api.Region
	if r, ok := utils.Servers[strings.ToLower(region)]; ok {
		reg= r
	}

	// check if the region was actually invalid
	if _, ok := utils.Servers[strings.ToLower(region)]; !ok {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem parsing league server")
		return
	}

	// forget the region for now
	client := golio.NewClient(os.Getenv("LEAGUE_API"), golio.WithRegion(reg))
	summoner, err := client.Riot.Summoner.GetByName(leagueName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting user from league api")
		return
	}

	_, err = client.Riot.Match.List(summoner.AccountID, 0, 1)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting match information")
		return
	}
}
