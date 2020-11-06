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

	var content string
	for _, play := range plays {
		content += fmt.Sprintf("[%s/%s/%s/%s]>>PP:%s>>%s\n",
			play.Count300, play.Count100, play.Count50, play.CountMiss, play.PP, play.Date)
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("Plays for %s", osuName),
			Value:  content,
			Inline: false,
		},
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Type = "rich"
	messageEmbed.Title = "Compare"
	messageEmbed.Color = 44504
	messageEmbed.Fields = fields

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
