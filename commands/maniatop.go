package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func ManiaTopHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	osuName, err := utils.GetOsuUsername(msg.Author.ID, args)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse osu name")
		return
	}

	topPlays, err := utils.GetModeTopPlays(osuName, "mania")
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	var fields []*discordgo.MessageEmbedField

	// we can use a loop since all the fields are similar in a sense
	for index := range topPlays {
		// load the beatmap so that we can get more information other than the ID
		beatmap, err := utils.GetOsuBeatmap(topPlays[index].BeatmapID)
		if err != nil {
			_, _ = session.ChannelMessage(msg.ChannelID,
				fmt.Sprintf("Error getting beatmap information on top score #%d", index+1))
			// if there was an error, still try to display rest of the top plays
			continue
		}

		formattedPP := strings.Split(topPlays[index].PP, ".")
		formattedValue := fmt.Sprintf("PP: %s, Score set: %s", formattedPP[0], topPlays[index].Date)

		// do all the needed bitwise calculations to get the mods; the error will never happen,
		// but handle it for good merit!
		mods, err := utils.GetMods(topPlays[index].EnabledMods)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		formattedTitle := fmt.Sprintf("%s + %s", beatmap.Title, mods)

		// finally add the new field to the fields array
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   formattedTitle,
			Value:  formattedValue,
			Inline: false})
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = fmt.Sprintf("osu! Mania top plays for %s", utils.UnFormatName(osuName))
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
