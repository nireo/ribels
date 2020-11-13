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

	// load beatmap to extract more information
	beatmap, err := utils.GetOsuBeatmap(currentMapID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not load beatmap information")
		return
	}

	// if there are no recent plays for that user
	if len(plays) == 0 {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("**Could not find plays by %s on %s[%s]**", osuName, beatmap.Title, beatmap.Version))
		return
	}

	userId := plays[0].UserID

	// precalculate all the PP values and IF FC values and the difficulty of maps with mods
	preCalculated, err := utils.PrecalculatePP(plays, currentMapID)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	var content string
	for index, play := range plays {
		mods, err := utils.GetMods(play.EnabledMods)
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse mod data")
			return
		}

		content += fmt.Sprintf("**%d.** `%s` **Score** [%.2f★]\n", index+1, mods, preCalculated[index].Diff)
		content += fmt.Sprintf("▸ %s ▸ **%.2fPP** *(%.2fpp for FC)*▸ %s%%\n",
			utils.RankEmojis[play.Rank], preCalculated[index].PlayPP, preCalculated[index].IfFCPP, play.CalculateAcc())
		content += fmt.Sprintf("▸ %s ▸ x%s/%s ▸ [%s/%s/%s/%s]\n",
			play.Score, play.MaxCombo, beatmap.MaxCombo, play.Count300, play.Count100, play.Count50, play.CountMiss)
		content += fmt.Sprintf("▸ Score set %s\n", play.Date)
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504
	messageEmbed.Description = content
	messageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf("https://b.ppy.sh/thumb/%sl.jpg", beatmap.BeatmapSetID),
	}
	messageEmbed.Author = &discordgo.MessageEmbedAuthor{
		IconURL: fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
		Name:    fmt.Sprintf("Top osu! Standard Plays for %s on %s[%s]", osuName, beatmap.Title, beatmap.Version),
		URL:     fmt.Sprintf("http://s.ppy.sh/a/%s", userId),
	}

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
