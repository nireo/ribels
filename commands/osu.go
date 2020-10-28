package commands

import (
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

	// create embed fields
	fields := []*discordgo.MessageEmbedField{
		{
			Name: "Playcount",
			Value: selectedUser.Playcount,
			Inline: false,
		},
		{
			Name: "Rank",
			Value: selectedUser.PPRank,
			Inline: false,
		},
		{
			Name: "Playtime",
			Value: selectedUser.SecondsPlayed,
			Inline: false,
		},
		{
			Name: "Level",
			Value: selectedUser.Level,
			Inline: false,
		},
		{
			Name: "Country",
			Value: selectedUser.Country,
			Inline: false,
		},
		{
			Name: "Accuracy",
			Value: selectedUser.Accuracy,
			Inline: false,
		},
	}

	// create the final embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = utils.UnFormatName(selectedUser.Username)
	messageEmbed.Fields = fields
	messageEmbed.Type = "rich"

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
