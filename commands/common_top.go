package commands

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func CommonTopCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string, mode string) {
	osuName, err := utils.GetOsuUsername(m.Author.ID, args)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting osu username")
		return
	}

	topPlays, err := utils.GetModeTopPlays(osuName, mode)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	userId := topPlays[0].UserID
	var content string

	for index, play := range topPlays {
		// load the beatmap so that we can get more information other than the ID
		beatmap, err := utils.GetOsuBeatmapMods(play.BeatmapID, play.EnabledMods)
		if err != nil {
			continue
		}

		// do all the needed bitwise calculations to get the mods; the error will never happen,
		// but handle it for good merit!
		mods, err := utils.GetMods(play.EnabledMods)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		ppFloat, err := strconv.ParseFloat(play.PP, 64)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		starFloat, err := strconv.ParseFloat(beatmap.Difficulty, 64)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		content += fmt.Sprintf("**%d. %s[%s] +%s** [%.2f★]\n",
			(index + 1), beatmap.Title, beatmap.Version, mods, starFloat)
		content += fmt.Sprintf("▸ %s ▸ **%.2f** ▸ %s%%\n",
			utils.RankEmojis[play.Rank], ppFloat, play.CalculateTopPlayAcc())
		content += fmt.Sprintf("▸ %s ▸ x%s/%s ▸ [%s/%s/%s/%s]\n",
			play.Score, play.MaxCombo, beatmap.MaxCombo, play.Count300, play.Count100, play.Count50, play.CountMiss)
		content += fmt.Sprintf("▸ Score Set %s\n\n", play.Date)
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504
	messageEmbed.Description = content
	messageEmbed.Footer = &discordgo.MessageEmbedFooter{
		Text: "On osu! Official Server",
	}
	messageEmbed.Author = &discordgo.MessageEmbedAuthor{
		IconURL: fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
		Name:    fmt.Sprintf("Top 3 osu! Standard plays for %s", osuName),
		URL:     fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
	}
	messageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
	}

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
