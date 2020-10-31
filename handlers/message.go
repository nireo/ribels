package handlers

import (
	"github.com/nireo/ribels/utils"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/commands"
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

	// check if logging is enabled and also check if the message starts the command flag
	// so that we don't log unrelated messages
	if logging && msg.Content[0] == '$' {
		log.Printf("%s : %s", msg.Author.ID, msg.Content)
	}

	// For each command create a goroutine, so basically concurrently execute all commands!
	switch args[0] {
	case utils.FormatCommand("set"):
		go commands.SetCommandHandler(session, msg, args)
	case utils.FormatCommand("recent"):
		go commands.RecentCommandHandler(session, msg, args)
	case utils.FormatCommand("osu"):
		go commands.OsuCommandHandler(session, msg, args)
	case utils.FormatCommand("top"):
		go commands.TopCommandHandler(session, msg, args)
	case utils.FormatCommand("help"):
		go commands.HelpCommandHandler(session, msg)
	case utils.FormatCommand("github"):
		go commands.GithubCommandHandler(session, msg)
	case utils.FormatCommand("maniatop"):
		go commands.ManiaTopHandler(session, msg, args)
	case utils.FormatCommand("taikotop"):
		go commands.TaikoTopCommandHandler(session, msg, args)
	case utils.FormatCommand("ctbtop"):
		go commands.CTBCommandHandler(session, msg, args)
	case utils.FormatCommand("map"):
		go commands.MapCommandHandler(session, msg, args)
	case utils.FormatCommand("set-lol"):
		go commands.SetLeagueCommandHandler(session, msg, args)
	case utils.FormatCommand("lol"):
		go commands.LeagueProfileCommandHandler(session, msg, args)
	case utils.FormatCommand("servers"):
		go commands.ServersCommandHandler(session, msg)
	case utils.FormatCommand("lol-remove"):
		go commands.RemoveLolCommandHandler(session, msg)
	case utils.FormatCommand("lol-curr"):
		go commands.CurrentLeagueGameCommand(session, msg, args)
	case utils.FormatCommand("osu-remove"):
		go commands.RemoveOsuCommandHandler(session, msg)
	// if we can't find a matching command just close the handler
	default:
		return
	}
}
