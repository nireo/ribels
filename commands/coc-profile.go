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

	_, _ = session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Town hall level %d", player.TownHallLevel))
}
