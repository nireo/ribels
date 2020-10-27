package commands

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

func CurrentLeagueGameCommand(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := ""
	sendMessage, err := session.ChannelMessageSend(msg.ChannelID, "Searching...")
	client := utils.NewRiotClient("euw1", os.Getenv("LEAGUE_API"), 10)
	summonerName := strings.Join(args[1:], " ")
	match, err := client.GetLiveMatchBySummonerName(&summonerName)

	if err == nil {
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCaption(true, fmt.Sprintf("Current match: %session", summonerName))
		table.SetHeader([]string{"Team", "Champion", "Summoner", "Solo"})

		var bt, rt []*discordgo.MessageEmbedField

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

		for i := range match {
			row := []*discordgo.MessageEmbedField{
				{
					Value:  match[i].Champion,
					Inline: true,
				},
				{
					Value:  match[i].SummonerName,
					Inline: true,
				},
				{
					Value:  match[i].Solo,
					Inline: true,
				},
			}

			if match[i].Team == "BLUE" {
				bt = append(bt, row...)
			} else if match[i].Team == "RED" {
				rt = append(rt, row...)
			}
		}

		var d [][]string
		for i := range match {
			d = append(d, []string{
				match[i].Team,
				match[i].Champion,
				match[i].SummonerName,
				match[i].Solo},
			)
		}

		table.AppendBulk(d)
		table.Render()
		content = "```" + buf.String() + "```"
	}

	if content == "" {
		content = "ERROR: Summoner most likely not in a game"
	}
	_, _ = session.ChannelMessageEdit(sendMessage.ChannelID, sendMessage.ID, content)
}
