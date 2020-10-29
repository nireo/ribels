package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
	"strconv"
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

	validRegion, err := utils.CheckValidRegion(region)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Error parsing region")
		return
	}

	client := utils.NewRiotClient(validRegion, os.Getenv("LEAGUE_API"))
	summoner, err := client.GetSummonerWithName(leagueName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not find user")
		return
	}

	ranks, err := client.GetSummonerRankWithID(summoner.ID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not get rank")
		return
	}

	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Username",
		Value:  summoner.Name,
		Inline: false,
	})

	levelString := strconv.Itoa(summoner.SummonerLevel)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Level",
		Value:  levelString,
		Inline: false,
	})

	for _, rank := range ranks {
		var name string
		if rank.QueueType == "RANKED_SOLO_5x5" {
			name = "Ranked Solo/Duo"
		} else {
			name = "Ranked Flex"
		}

		content := fmt.Sprintf("%s %s %d | W/L [%d/%d]",
			rank.Tier, rank.Rank, rank.LeaguePoints, rank.Wins, rank.Losses)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  content,
			Inline: false,
		})
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "Profile information"
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
