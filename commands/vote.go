package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func VoteCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	question := strings.Join(args[1:], " ")
	if question == "" {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "Cannot vote on a empty question!")
		return
	}

	m, _ := session.ChannelMessageSend(msg.ChannelID, question)
	session.MessageReactionAdd(msg.ChannelID, m.ID, "👍")
	session.MessageReactionAdd(msg.ChannelID, m.ID, "👎")

	time.Sleep(time.Second * 5)

	for reaction := range m.Reactions {
		fmt.Println(reaction)
	}
}
