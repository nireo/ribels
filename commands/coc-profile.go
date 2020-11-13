package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func CocProfileCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	cocTag := args[1]
	player, err := utils.GetCOCPlayerData(cocTag)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	var content string
	content += fmt.Sprintf("**▸ TH Level:** %d\n", player.TownHallLevel)
	content += fmt.Sprintf("**▸ League :** %s (%d)\n", player.League.Name, player.Trophies)
	content += fmt.Sprintf("**▸ Clan :** %s\n", player.Clan.Name)

	content += "\n** Heroes **\n"
	for _, hero := range player.Heroes {
		if hero.Village == "home" {
			content += fmt.Sprintf("**▸ %s :** %d/%d\n", hero.Name, hero.Level, hero.MaxLevel)
		}
	}

	content += "\n** Troops **\n"
	for _, troop := range player.Troop {
		if troop.Village == "home" {
			content += fmt.Sprintf("**▸ %s :** %d/%d\n", troop.Name, troop.Level, troop.MaxLevel)
		}
	}

	content += "\n** Spells **\n"
	for _, spell := range player.Spells {
		if spell.Village == "home" {
			content += fmt.Sprintf("**▸ %s :** %d/%d\n", spell.Name, spell.Level, spell.MaxLevel)
		}
	}

	var messageEmbed discordgo.MessageEmbed
	messageEmbed.Description = content
	messageEmbed.Type = "rich"
	messageEmbed.Color = 44504

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
