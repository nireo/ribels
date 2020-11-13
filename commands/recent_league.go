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
	match := matches.Matches[0]

	matchInfo, err := client.GetSingleMatch(strconv.Itoa(match.GameID))
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	champion := client.Champions.GetChampionWithKey(strconv.Itoa(match.Champion))
	content += fmt.Sprintf("\n**▸ Champion:** %s\n", champion.Name)
	content += fmt.Sprintf("**▸ Game Duration:** %d minutes\n", matchInfo.GameDuration/60)

	var participantID int
	for _, participant := range matchInfo.ParticipantIdentities {
		if participant.Player.SummonerID == summoner.ID {
			participantID = participant.ParticipantID
		}
	}

	var par utils.MatchParticipant
	for _, participant := range matchInfo.Participants {
		if participantID == participant.ParticipantId {
			par = participant
		}
	}

	content += fmt.Sprintf("**▸ Stats:** %d/%d/%d\n",
		par.Stats.Kills, par.Stats.Deaths, par.Stats.Assists)

	t := time.Unix(match.Timestamp/1000, 0)

	strDate := t.Format(time.UnixDate)
	content += fmt.Sprintf("**▸ Played:** %s\n", strDate)

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
