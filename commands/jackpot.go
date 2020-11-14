package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
	"strconv"
	"time"
)

func JackpotCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	switch args[1] {
	case "start":
		if utils.GameRunning {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "A game is already running wait for the results!")
			return
		}

		if len(args) <= 2 {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "You need to provide 2 arguments")
			return
		}

		// check that the user has added a starting bid
		wager, err := strconv.Atoi(args[2])
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "You need to add a starting wager!")
			return
		}

		utils.StartGame()
		if err := utils.AddPlayer(msg.Author.ID, msg.Author.Username, int64(wager)); err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Error creating the jackpot")
		}

		// create a sleep timer, and this is the same command that will display the results of the jackpot
		time.Sleep(time.Second * 15)

		_, _ = session.ChannelMessageSend(msg.ChannelID, "The game has concluded")
		utils.ClearGame()
	case "join":
		if !utils.GameRunning {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "There is no game currently")
			return
		}

		if len(args) <= 2 {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "You need to provide 2 arguments")
			return
		}

		// check that the user has added a starting bid
		wager, err := strconv.Atoi(args[2])
		if err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "You need to add a starting wager!")
			return
		}

		if err := utils.AddPlayer(msg.Author.ID, msg.Author.Username, int64(wager)); err != nil {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Error creating the jackpot")
			return
		}

		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Your bet of `%d` has been added! Current standings:\n%s", wager, utils.PrintPlayers()))
	case "players":
		if utils.GameRunning {
			players := utils.PrintPlayers()
			_, _ = session.ChannelMessageSend(msg.ChannelID, players)
			return
		}
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Game is not currently running")
		return
	default:
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			"Jackpot command not found, supported commands `players join start`")
		return
	}
}
