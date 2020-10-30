package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
	"sort"
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

		winRate := float64(rank.Wins)/float64(rank.Wins+rank.Losses)

		content := fmt.Sprintf("%s %s %d | W/L [%d/%d] %0.f%%",
			rank.Tier, rank.Rank, rank.LeaguePoints, rank.Wins, rank.Losses, winRate*100.0)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  content,
			Inline: false,
		})
	}

	// get the most played champs
	masteries, err := client.ListsSummonerMasteries(summoner.ID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem getting champion masteries")
		return
	}

	sort.SliceStable(masteries, func(a, b int) bool {
		return masteries[a].ChampionPoints > masteries[b].ChampionPoints
	})

	fieldValue := ""
	for _, champion := range masteries[:5] {
		championKey := strconv.Itoa(champion.ChampionID)
		champ := client.Champions.GetChampionWithKey(championKey)
		if champ != nil {
			fieldValue += fmt.Sprintf("%s (%d)\n", champ.Name, champion.ChampionPoints)
		}
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: "Most played champions",
		Value: fieldValue,
		Inline: false,
	})

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "Profile information"
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
