package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

// This command gives information about a certain user,
// either the user from an argument or a user from the database
func OsuCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	// check if a user argument is provided, otherwise load user from database
	var osuName string
	if len(args) > 1 {
		osuName = utils.FormatName(args[1:])
	} else {
		user, err := utils.CheckIfSet(msg.Author.ID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID,
				"Your osu! profile is not set. To do this type $set osu_name")
		}

		osuName = user.OsuName
	}

	// The osu api gives every single request as an array so we just need to extract the first element
	selectedUser, err := utils.GetUserFromOSU(osuName)
	if err != nil {
		session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	// create embed fields
	var fields []*discordgo.MessageEmbedField
	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Playcount", Value: selectedUser.Playcount, Inline: false})

	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Rank", Value: selectedUser.PPRank, Inline: false})

	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Playtime", Value: selectedUser.SecondsPlayed, Inline: false})

	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Level", Value: selectedUser.Level, Inline: false})

	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Country", Value: selectedUser.Country, Inline: false})

	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Accuracy", Value: selectedUser.Accuracy, Inline: false})

	// create the final embed
	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = utils.UnFormatName(selectedUser.Username)
	messageEmbed.Fields = fields
	messageEmbed.Type = "rich"

	session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
