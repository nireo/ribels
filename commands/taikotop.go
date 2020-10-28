package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func TaikoTopCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	osuName, err := utils.GetOsuUsername(msg.Author.ID, args)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Error getting osu name")
		return
	}

	topPlays, err := utils.GetModeTopPlays(osuName, "mania")
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	var fields []*discordgo.MessageEmbedField

	// we can use a loop since all the fields are similar in a sense
	for _, play := range topPlays {
		// load the beatmap so that we can get more information other than the ID
		beatmap, err := utils.GetOsuBeatmap(play.BeatmapID)
		if err != nil {
			continue
		}

		formattedPP := strings.Split(play.PP, ".")
		formattedValue := fmt.Sprintf("PP: %s, Score set: %s", formattedPP[0], play.Date)

		// do all the needed bitwise calculations to get the mods; the error will never happen,
		// but handle it for good merit!
		mods, err := utils.GetMods(play.EnabledMods)
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
	messageEmbed.Title = fmt.Sprintf("osu! Taiko top plays for %s", utils.UnFormatName(osuName))
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
