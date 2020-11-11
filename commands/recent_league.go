package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func RecentLeagueCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var region string
	var user utils.LeagueUser
	var username string

	db := utils.GetDatabase()
	if len(args) != 3 {
		if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting user from database!")
			return
		}

		region = user.Region
		username = user.Username
	} else {
		username = args[1]
		region = args[2]
	}

	validRegion, err := utils.CheckValidRegion(region)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Error parsing region")
		return
	}

	client := utils.NewRiotClient(validRegion, os.Getenv("LEAGUE_API"))
	summoner, err := client.GetSummonerWithName(username)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Could not find summoner `%s` on `%s`", username, validRegion))
		return
	}

	matches, err := client.GetListOfMatches(summoner.AccountID, 1, 10)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Error `%s` while getting match data.", err.Error()))
		return
	}

	var content string
	for index, match := range matches.Matches[:3] {

		champion := client.Champions.GetChampionWithKey(strconv.Itoa(match.Champion))
		content += fmt.Sprintf("\n**%d. %s**\n", (index + 1), champion.Name)

		betterLane, ok := utils.Roles[match.Role]
		if !ok {
			content += "**▸ Lane:** Non specified\n"
		} else {
			content += fmt.Sprintf("**▸ Lane:** %s\n", betterLane)
		}

		t := time.Unix(match.Timestamp/1000, 0)
		fmt.Println(match.Timestamp)

		strDate := t.Format(time.UnixDate)
		content += fmt.Sprintf("**▸ Played:** %s\n", strDate)
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("Most recent games for %s", summoner.Name),
			Value:  content,
			Inline: false,
		},
	}
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
