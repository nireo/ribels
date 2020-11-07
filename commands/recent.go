package commands

import (
	"fmt"
	"strconv"

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

	mods, err := utils.GetMods(recentPlay.EnabledMods)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Error parsing mods")
		return
	}

	floatDifficulty, err := strconv.ParseFloat(beatmap.Difficulty, 64)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse star rating")
		return
	}

	var content string
	content += fmt.Sprintf("▸ %s ▸ PP ▸ %s%%\n", utils.RankEmojis[recentPlay.Rank], recentPlay.CalculateAcc())
	content += fmt.Sprintf("▸ %s ▸ x%s/%s ▸ [%s/%s/%s/%s]\n",
		recentPlay.Score, recentPlay.MaxCombo, beatmap.MaxCombo, recentPlay.Count300,
		recentPlay.Count100, recentPlay.Count50, recentPlay.CountMiss)

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("%s[%s] + %s[%.2f★]", beatmap.Title, beatmap.Version, mods, floatDifficulty),
			Value:  content,
			Inline: false,
		},
	}

	footer := *&discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Score set %s", recentPlay.Date),
	}

	// create the actual embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields
	messageEmbed.Color = 44504
	messageEmbed.Footer = &footer

	// set this recent map as the current map, so that others can compare their scores
	utils.SetCurrentMap(recentPlay.BeatmapID)

	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("**Most recent osu!Standard play for %s**", utils.UnFormatName(osuName)))

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
