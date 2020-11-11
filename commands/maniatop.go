package commands

import (
	"fmt"
	"strconv"

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

	var content string

	// we can use a loop since all the fields are similar in a sense
	for index, play := range topPlays {
		// load the beatmap so that we can get more information other than the ID
		beatmap, err := utils.GetOsuBeatmap(topPlays[index].BeatmapID)
		if err != nil {
			_, _ = session.ChannelMessage(msg.ChannelID,
				fmt.Sprintf("Error getting beatmap information on top score #%d", index+1))
			// if there was an error, still try to display rest of the top plays
			continue
		}

		// do all the needed bitwise calculations to get the mods; the error will never happen,
		// but handle it for good merit!
		mods, err := utils.GetMods(topPlays[index].EnabledMods)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		ppFloat, err := strconv.ParseFloat(play.PP, 64)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		starFloat, err := strconv.ParseFloat(beatmap.Difficulty, 64)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
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

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("Top 3 osu! Mania Plays for %s", osuName),
			Value:  content,
			Inline: false,
		},
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Fields = fields
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
