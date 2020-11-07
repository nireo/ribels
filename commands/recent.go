package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func RecentCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	osuName, err := utils.GetOsuUsername(msg.Author.ID, args)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse osu user")
		return
	}

	// get the most recent play from user
	recentPlay, err := utils.GetRecentPlay(osuName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	// get the beatmap, so that we can use it's name and other data related to it
	beatmap, err := utils.GetOsuBeatmap(recentPlay.BeatmapID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	// make fields for the rich embed message
	formattedCounts := fmt.Sprintf("[%s/%s/%s/%s]",
		recentPlay.Count300, recentPlay.Count100, recentPlay.Count50, recentPlay.CountMiss)
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Counts",
			Value:  formattedCounts,
			Inline: true,
		},
		{
			Name:   "Score",
			Value:  recentPlay.Score,
			Inline: true,
		},
		{
			Name:   "Combo",
			Value:  fmt.Sprintf("[%sx/%sx]", recentPlay.MaxCombo, beatmap.MaxCombo),
			Inline: true,
		},
		{
			Name:   "Date",
			Value:  recentPlay.Date,
			Inline: true,
		},
		{
			Name:   "Difficulty",
			Value:  beatmap.Difficulty,
			Inline: true,
		},
		{
			Name:   "Rank",
			Value:  utils.RankEmojis[recentPlay.Rank],
			Inline: true,
		},
	}

	// create the actual embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Title = fmt.Sprintf("%s[%s]", beatmap.Title, beatmap.Version)
	messageEmbed.Fields = fields
	messageEmbed.Color = 44504

	// set this recent map as the current map, so that others can compare their scores
	utils.SetCurrentMap(recentPlay.BeatmapID)

	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("Most recent osu!Standard play for %s", utils.UnFormatName(osuName)))

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
