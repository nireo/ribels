package handlers

import (
	"log"
	"strings"

	"github.com/nireo/ribels/commands"
	"github.com/nireo/ribels/utils"

	"github.com/bwmarrin/discordgo"
)

var logging bool

// Use this approach, since we don't want to load ENV variables for each message!
func SetLogging(status bool) {
	logging = status
}

func MessageHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// ignore all messages created by the bot.
	if msg.Author.ID == session.State.User.ID {
		return
	}
	// tokenize the input
	args := strings.Split(msg.Content, " ")

	levenshteinDis := utils.ComputeDistance(strings.ReplaceAll(msg.Content, " ", ""), "dead")
	if levenshteinDis < 4 {
		session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		_, _ = session.ChannelMessageSend(msg.ChannelID, "https://tenor.com/view/dead-chat-passione-admin-passione-jojolion-gif-19211422")
		return
	}

	// check if logging is enabled and also check if the message starts the command flag
	// so that we don't log unrelated messages
	if logging && msg.Content[0] == '$' {
		log.Printf("%s : %s", msg.Author.ID, msg.Content)
	}

	// check for the commands with arguments
	if command, ok := commands.CommandsWithArgs[strings.ToLower(args[0])]; ok {
		go command(session, msg, args)
		// stop needless checking after executing command
		return
	}

	if command, ok := commands.CommandsWOArgs[strings.ToLower(args[0])]; ok {
		go command(session, msg)
		return
	}
}
