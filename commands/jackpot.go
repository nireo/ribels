package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"time"
)

func JackpotCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	switch args[1] {
	case "start": {
		if utils.GameRunning {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "A game is already running wait for the results!")
			return
		} else {
			// create a sleep timer, and this is the same command that will display the results of the jackpot
			utils.GameRunning = true
			time.Sleep(time.Second*10)

			_, _ = session.ChannelMessageSend(msg.ChannelID, "The game has concluded")
			utils.ClearGame()
		}
	}
	case "players":
		if utils.GameRunning {
			players := utils.PrintPlayers()
			_, _ = session.ChannelMessageSend(msg.ChannelID, players)
			return
		} else {
			_, _ = session.ChannelMessageSend(msg.ChannelID,
				"Game is not currently running")
			return
		}
	default:
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Jackpot command not found, supported commands `players join start`")
		return
	}
}
