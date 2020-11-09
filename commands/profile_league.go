package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	levelString := strconv.Itoa(summoner.SummonerLevel)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Level",
		Value:  levelString,
		Inline: true,
	})

	t := time.Unix(summoner.Revisiondate, 0)
	strDate := t.Format(time.Stamp)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Last played",
		Value:  strDate,
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "\u200B",
		Value:  "\u200B",
		Inline: false,
	})

	for _, rank := range ranks {
		var name string
		if rank.QueueType == "RANKED_SOLO_5x5" {
			name = "*Ranked Solo/Duo*"
		} else {
			name = "*Ranked Flex*"
		}

		winRate := float64(rank.Wins) / float64(rank.Wins+rank.Losses)

		formattedDivision := fmt.Sprintf(strings.Title(strings.ToLower(rank.Tier)))

		content := fmt.Sprintf("**%s %s** \n **%d LP**  %dW %dL \n Win Ratio %0.f%%",
			formattedDivision, rank.Rank, rank.LeaguePoints, rank.Wins, rank.Losses, winRate*100.0)

		if rank.MiniSeries.Progress != "" {
			content += fmt.Sprintf("\nPromos: %dW %dL", rank.MiniSeries.Wins, rank.MiniSeries.Losses)
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  content,
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "\u200B",
		Value:  "\u200B",
		Inline: false,
	})

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
		Name:   "Most played champions",
		Value:  fieldValue,
		Inline: false,
	})

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "Summoner " + summoner.Name
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
