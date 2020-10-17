package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func TopCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var osuName string
	if len(args) > 1 {
		osuName = utils.FormatName(args[1:])
	} else {
		user, err := utils.CheckIfSet(msg.Author.ID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Not set in database")
			return
		}

		osuName = user.OsuName
	}

	topPlays, err := utils.GetUserTopplaysFromOSU(osuName)
	if err != nil {
		session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	var fields []*discordgo.MessageEmbedField
	for index := range topPlays {
		beatmap, err := utils.GetOsuBeatmap(topPlays[index].BeatmapID)
		if err != nil {
			session.ChannelMessage(msg.ChannelID,
				fmt.Sprintf("Error getting beatmap information on top score #%d", index+1))
			return
		}

		formattedPP := strings.Split(topPlays[index].PP, ".")
		formattedValue := fmt.Sprintf("PP: %s, Score set: %s", formattedPP[0], topPlays[index].Date)

		mods, err := utils.GetMods(topPlays[index].EnabledMods)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		formattedTitle := fmt.Sprintf("%s + %s", beatmap.Title, mods)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   formattedTitle,
			Value:  formattedValue,
			Inline: false})
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = fmt.Sprintf("osu! Standard top plays for %s", utils.UnFormatName(osuName))
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields

	session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
