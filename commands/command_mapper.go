package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

type WithArgs func(*discordgo.Session, *discordgo.MessageCreate, []string)
type WOArgs func(*discordgo.Session, *discordgo.MessageCreate)

var CommandsWithArgs map[string]WithArgs
var CommandsWOArgs map[string]WOArgs

func InitCommandsMap() {
	CommandsWithArgs = map[string]WithArgs {
		utils.FormatCommand("set"): SetCommandHandler,
		utils.FormatCommand("recent"): RecentCommandHandler,
		utils.FormatCommand("osu"): OsuCommandHandler,
		utils.FormatCommand("maniatop"): ManiaTopHandler,
		utils.FormatCommand("ctbtop"): CTBCommandHandler,
		utils.FormatCommand("taikotop"): TaikoTopCommandHandler,
		utils.FormatCommand("map"): MapCommandHandler,
		utils.FormatCommand("set-lol"): SetLeagueCommandHandler,
		utils.FormatCommand("lol-curr"): CurrentLeagueGameCommand,
		utils.FormatCommand("lol"): LeagueProfileCommandHandler,
	}

	CommandsWOArgs = map[string]WOArgs {
		utils.FormatCommand("help"): HelpCommandHandler,
		utils.FormatCommand("servers"): ServersCommandHandler,
		utils.FormatCommand("lol-remove"): RemoveLolCommandHandler,
		utils.FormatCommand("osu-remove"): RemoveLolCommandHandler,
	}
}
