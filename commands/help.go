package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Doesn't need arguments!
func HelpCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	var content string
	content += "** osu! Commands:**\n `osu` `set` `top` `maniatop` `taikotop` `ctbtop` `rs` `c` `osu-rm` `map`\n"
	content += "** LoL commands:**\n `rs-lol` `set-lol` `lol-curr` `lol` `servers` `lol-rm`\n"
	content += "** Gamling:**\n `jackpot` `reset-balance` `balance` `coinflip` \n"
	content += "** Misc:**\n `help` `coc` "

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Title = "List of all ribels commands"
	messageEmbed.Type = "rich"
	messageEmbed.Description = content
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
