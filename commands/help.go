package commands

import "github.com/bwmarrin/discordgo"

// Doesn't need arguments!
func HelpCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "$set",
		Value:  "Link a osu! username to your discord id. Usage: $set osu_name",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "$osu",
		Value:  "List some information about a given user, if no username argument is provided the linked user will be used!",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "$top",
		Value:  "List all the top plays of the given user, if no username argument is provided the linked user will be used!",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "$recent",
		Value:  "List the most recent play of the given user, if no username argument is provided the linked user will be used!",
		Inline: false,
	})

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "ribels commands help"
	messageEmbed.Fields = fields
	messageEmbed.Type = "rich"
	messageEmbed.Description = "Every command in ribels, and the usage of the those commands."

	session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
