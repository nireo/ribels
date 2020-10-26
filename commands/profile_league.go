package commands

import (
	"fmt"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
	"strconv"
	"strings"
)

func LeagueProfileCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
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

	leagues, err := client.Riot.League.ListBySummoner(summoner.ID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem loading user league")
		return
	}

	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name: "Username",
		Value: summoner.Name,
		Inline: false,
	})

	levelString := strconv.Itoa(summoner.SummonerLevel)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name: "Level",
		Value: levelString,
		Inline: false,
	})

	for _, league := range leagues {
		leagueName := ""
		if league.QueueType == "RANKED_SOLO_5x5" {
			leagueName = "Ranked Solo/Duo"
		} else {
			leagueName = "Ranked Flex"
		}

		// format the values nicely
		content := fmt.Sprintf("%s %s %d | W/L [%d/%d]",
			league.Tier, league.Rank, league.LeaguePoints,
			league.Wins, league.Losses)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: leagueName,
			Value: content,
			Inline: false,
		})
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "Profile information"
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
