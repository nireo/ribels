package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

// This command gives information about a certain user,
// either the user from an argument or a user from the database
func OsuCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	// check if a user argument is provided, otherwise load user from database
	osuName, err := utils.GetOsuUsername(msg.Author.ID, args)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse osu name")
		return
	}

	// The osu api gives every single request as an array so we just need to extract the first element
	selectedUser, err := utils.GetUserFromOSU(osuName)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	var content string

	content += fmt.Sprintf("**▸ Official Rank:** #%s (%s#%s)\n",
		selectedUser.PPRank, selectedUser.Country, selectedUser.PPCountryRank)

	splittedLevel := strings.Split(selectedUser.Level, ".")

	formattedProgress := fmt.Sprintf("%s.%s", splittedLevel[1][:2], splittedLevel[1][2:])

	content += fmt.Sprintf("**▸ Level:** %s (%s%%)\n", splittedLevel[0], formattedProgress)
	content += fmt.Sprintf("**▸ Total PP:** %s\n", selectedUser.RawPP)

	floatAcc, err := strconv.ParseFloat(selectedUser.Accuracy, 64)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Could not parse hit accuracy")
	}
	content += fmt.Sprintf("**▸ Hit Accuracy:** %.2f%%\n", floatAcc)
	content += fmt.Sprintf("**▸ Playcount:** %s\n", selectedUser.Playcount)

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("osu! Standard Profile for %s", selectedUser.Username),
			Value:  content,
			Inline: false,
		},
	}

	// create the final embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Fields = fields
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504
	messageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf("http://s.ppy.sh/a/%s", selectedUser.UserID),
	}
	messageEmbed.Footer = &discordgo.MessageEmbedFooter{
		Text: "On osu! Official Server",
	}

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
