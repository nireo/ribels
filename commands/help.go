package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

// Doesn't need arguments!
func HelpCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("set"),
		Value:  "Link a osu! username to your discord id. Usage: $set osu_name",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("osu"),
		Value:  "List some information about a given user",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("top"),
		Value:  "List all the standard top plays of the given user",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("maniatop"),
		Value:  "List all the mania top plays of the given user",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("taikotop"),
		Value:  "List all the taiko top plays of the given user",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("ctbtop"),
		Value:  "List all the top plays of the given user",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("recent"),
		Value:  "List the most recent osu!standard play of the given user",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: utils.FormatCommand("set-lol"),
		Value: "Link a league of legends profile to your discord id",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: utils.FormatCommand("lol"),
		Value: "Show information about a league of legends player, [arguments: username server]",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: utils.FormatCommand("lol-curr"),
		Value: "Show game information about a current game",
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   utils.FormatCommand("rs-lol"),
		Value:  "Show the last game a summoner has played",
		Inline: false,
	})

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "ribels commands help"
	messageEmbed.Fields = fields
	messageEmbed.Type = "rich"
	messageEmbed.Description = "Every command in ribels, and the usage of the those commands."
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
