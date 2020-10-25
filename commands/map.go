package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func MapCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	// the user can give an url which has the beatmap id
	tokenizedURL := strings.Split(args[1], "/")

	// the beatmap id is stored in the last item
	beatmapID := tokenizedURL[len(tokenizedURL)-1]

	// sometimes the beatmapID has the difficulty included, so check for that
	index := strings.Index(beatmapID, "#")
	if index != -1 {
		beatmapID = beatmapID[:index]
	}

	// find the osu beatmap
	beatmap, err := utils.GetOsuBeatmap(beatmapID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Problem retrieving osu map info")
		return
	}

	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Title",
		Value:  beatmap.Title,
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "BPM",
		Value:  beatmap.BPM,
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Length",
		Value:  beatmap.TotalLength,
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Max combo",
		Value:  beatmap.MaxCombo,
		Inline: false,
	})

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "Map information"
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
