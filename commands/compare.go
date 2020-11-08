package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func CompareCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	osuName, err := utils.GetOsuUsername(msg.Author.ID, args)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse osu user.")
		return
	}

	plays, err := utils.GetScoresForCurrentMap(osuName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not get recent plays for user.")
		return
	}

	// load the beatmap, so that we can display it better
	currentMapID := utils.GetCurrentMap()
	if currentMapID == "" {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "**No map in discussion**")
		return
	}

	beatmap, err := utils.GetOsuBeatmap(currentMapID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not load beatmap information")
		return
	}

	if len(plays) == 0 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("**Could not plays by %s on %s**", osuName, beatmap.Title))
		return
	}

	var content string
	for index, play := range plays {
		mods, err := utils.GetMods(play.EnabledMods)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse mod data")
			return
		}

		content += fmt.Sprintf("**%d.** `%s` **Score** [%s★]\n", (index + 1), mods, beatmap.Difficulty)
		content += fmt.Sprintf("▸ %s ▸ x%s/%s ▸ [%s/%s/%s/%s]\n",
			play.Score, play.MaxCombo, beatmap.MaxCombo, play.Count300, play.Count100, play.Count50, play.Count50)
		content += fmt.Sprintf("▸ Score set %s\n\n", play.Date)
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("Plays by %s on %s[%s]", osuName, beatmap.Title, beatmap.Version),
			Value:  content,
			Inline: false,
		},
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504
	messageEmbed.Fields = fields
	messageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf("https://b.ppy.sh/thumb/%sl.jpg", beatmap.BeatmapSetID),
	}

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
