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

	userId := recentPlay.UserID

	// get the beatmap, so that we can use it's name and other data related to it
	beatmap, err := utils.GetOsuBeatmapMods(recentPlay.BeatmapID, recentPlay.EnabledMods)
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

	calculatedPP, err := recentPlay.CalculatePP()
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse osu data")
		return
	}

	var content string
	content += fmt.Sprintf("▸ %s ▸ **%.2fPP** *(%.2fpp for FC)*▸ %s%%\n",
		utils.RankEmojis[recentPlay.Rank], calculatedPP.PlayPP, calculatedPP.IfFCPP, recentPlay.CalculateAcc())

	content += fmt.Sprintf("▸ %s ▸ x%s/%s ▸ [%s/%s/%s/%s]\n",
		recentPlay.Score, recentPlay.MaxCombo, beatmap.MaxCombo, recentPlay.Count300,
		recentPlay.Count100, recentPlay.Count50, recentPlay.CountMiss)

	// create the actual embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504
	messageEmbed.Description = content
	messageEmbed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Score set %s | https://osu.ppy.sh/b/%s", recentPlay.Date, recentPlay.BeatmapID),
	}
	messageEmbed.Author = &discordgo.MessageEmbedAuthor{
		IconURL: fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
		Name:    fmt.Sprintf("%s[%s] + %s[%.2f★]", beatmap.Title, beatmap.Version, mods, floatDifficulty),
		URL:     fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
	}

	messageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf("https://b.ppy.sh/thumb/%sl.jpg", beatmap.BeatmapSetID),
	}

	// set this recent map as the current map, so that others can compare their scores
	utils.SetCurrentMap(recentPlay.BeatmapID)

	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("**Most recent osu!Standard play for %s**", utils.UnFormatName(osuName)))

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
