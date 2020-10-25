package commands

import (
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
	"strconv"
)

func LeagueProfileCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var leagueName string
	var user utils.LeagueUser

	db := utils.GetDatabase()
	if len(args) != 3 {
		if err := db.Where(&utils.LeagueUser{DiscordID: msg.Author.ID}).First(&user).Error; err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting user from database!")
			return
		}

		leagueName = user.Username
	} else {
		leagueName = args[1]
	}

	// forget the region for now
	client := golio.NewClient(os.Getenv("LEAGUE_API"), golio.WithRegion(api.RegionEuropeWest))
	summoner, err := client.Riot.Summoner.GetByName(leagueName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting user from league api")
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

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "Profile information"
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
