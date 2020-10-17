package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/commands"
	"github.com/nireo/ribels/utils"
)

func MessageHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// tokenize the input
	args := strings.Split(msg.Content, " ")

	switch args[0] {
	case "$set":
		go commands.SetCommandHandler(session, msg, args)
	case "$osu":
		go commands.OsuCommandHandler(session, msg, args)
	case "$top":
		go commands.TopCommandHandler(session, msg, args)
	default:
		return
	}

	if args[0] == "$recent" {
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

		// get the most recent play from user
		recentPlay, err := utils.GetRecentPlay(osuName)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

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

		beatmap, err := utils.GetOsuBeatmap(recentPlay.BeatmapID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		var messageEmbed discordgo.MessageEmbed
		messageEmbed.Title = beatmap.Title
		messageEmbed.Fields = fields

		session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Most recent osu!Standard play for %s", utils.UnFormatName(osuName)))

		session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
	}
}
