package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func ExecCommandHandler(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) <= 2 {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "You need to add more arguments!")
		return
	}

	code := strings.Join(args[2:], "")
	executionInformation, err := utils.ExecuteCodeRequest(args[1], code)
	if err != nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("took %s seconds\noutput:\n```%s```", executionInformation.CPUTime, executionInformation.Output))
}
