package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func RecentCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var osuName string

	// check if the user has provided an argument, otherwise load the osuname from the database
	if len(args) > 1 {
		osuName = utils.FormatName(args[1:])
	} else {
		user, err := utils.CheckIfSet(msg.Author.ID)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Not set in database")
			return
		}

		osuName = user.OsuName
	}

	// get the most recent play from user
	recentPlay, err := utils.GetRecentPlay(osuName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	// make fields for the rich embed message
	var fields []*discordgo.MessageEmbedField
	formattedCounts := fmt.Sprintf("[%s/%s/%s/%s]",
		recentPlay.Count300, recentPlay.Count100, recentPlay.Count50, recentPlay.CountMiss)

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Counts",
		Value:  formattedCounts,
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Score",
		Value:  recentPlay.Score,
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Date",
		Value:  recentPlay.Date,
		Inline: false,
	})

	// get the beatmap, so that we can use it's name and other data related to it
	beatmap, err := utils.GetOsuBeatmap(recentPlay.BeatmapID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	// create the actual embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = beatmap.Title
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("Most recent osu!Standard play for %s", utils.UnFormatName(osuName)))

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
