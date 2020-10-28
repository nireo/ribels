package commands

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"github.com/olekukonko/tablewriter"
	"os"
)

func CurrentLeagueGameCommand(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := ""
	sendMessage, err := session.ChannelMessageSend(msg.ChannelID, "Searching...")
	region, err := utils.CheckValidRegion(args[2])
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem with parsing server")
		return
	}

	client := utils.NewRiotClient(region, os.Getenv("LEAGUE_API"))
	summonerName := args[1]
	match, err := client.GetLiveMatchBySummonerName(summonerName)
	var bt, rt []*discordgo.MessageEmbedField
	if err == nil {
		// create a new writer for the ascii table
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)

		// table settings
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCaption(true, fmt.Sprintf("Current match: %session", summonerName))
		table.SetHeader([]string{"Team", "Champion", "Summoner", "Solo"})

		// create message embed fields for both teams
		headers := []string{"Champion", "Summoner", "Solo"}
		for i := range headers {
			bt = append(bt, &discordgo.MessageEmbedField{
				Name:   headers[i],
				Value:  "\u200b",
				Inline: true,
			})
			rt = append(rt, &discordgo.MessageEmbedField{
				Name:   headers[i],
				Value:  "\u200b",
				Inline: true,
			})
		}

		// for each player in the match add a message embed field
		for _, player := range match {
			row := []*discordgo.MessageEmbedField{
				{
					Value:  player.Champion,
					Inline: true,
				},
				{
					Value:  player.SummonerName,
					Inline: true,
				},
				{
					Value:  player.Solo,
					Inline: true,
				},
			}

			if player.Team == "BLUE" {
				bt = append(bt, row...)
			} else if player.Team == "RED" {
				rt = append(rt, row...)
			}
		}

		//create the table
		var tab [][]string
		for _, m := range match {
			// add row to table
			tab = append(tab, []string{
				m.Team,
				m.Champion,
				m.SummonerName,
				m.Solo},
			)
		}

		table.AppendBulk(tab)
		table.Render()

		// ``` are used this way discord makes the message format more flexible
		content = "```" + buf.String() + "```"
	}

	// if there was no table, an error has occured so we prompt the user
	if content == "" {
		content = "Summoner most likely not in a game"
	}

	// give out the table with the current game
	_, _ = session.ChannelMessageEdit(sendMessage.ChannelID, sendMessage.ID, content)
}
