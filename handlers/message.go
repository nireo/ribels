package handlers

import (
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
	// tokenize the input
	args := strings.Split(msg.Content, " ")

	// check if logging is enabled and also check if the message starts the command flag
	// so that we don't log unrelated messages
	if logging && msg.Content[0] == '$' {
		log.Printf("%s : %s", msg.Author.ID, msg.Content)
	}

	// For each command create a goroutine, so basically concurrently execute all commands!
	switch args[0] {
	case "$set":
		go commands.SetCommandHandler(session, msg, args)
	case "$recent":
		go commands.RecentCommandHandler(session, msg, args)
	case "$osu":
		go commands.OsuCommandHandler(session, msg, args)
	case "$top":
		go commands.TopCommandHandler(session, msg, args)
	case "$help":
		go commands.HelpCommandHandler(session, msg)
	case "$github":
		go commands.GithubCommandHandler(session, msg)
	case "$maniatop":
		go commands.ManiaTopHandler(session, msg, args)
	case "$taikotop":
		go commands.TaikoTopCommandHandler(session, msg, args)
	case "$ctbtop":
		go commands.CTBCommandHandler(session, msg, args)
	case "$map":
		go commands.MapCommandHandler(session, msg, args)
	case "$set-league":
		go commands.SetLeagueCommandHandler(session, msg, args)
	case "$lol-profile":
		go commands.LeagueProfileCommandHandler(session, msg, args)
	// if we can't find a matching command just close the handler
	default:
		return
	}
}
