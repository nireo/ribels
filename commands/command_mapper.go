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
	CommandsWithArgs = map[string]WithArgs{
		utils.FormatCommand("set"):      SetCommandHandler,
		utils.FormatCommand("top"):      TopCommandHandler,
		utils.FormatCommand("rs"):       RecentCommandHandler,
		utils.FormatCommand("osu"):      OsuCommandHandler,
		utils.FormatCommand("maniatop"): ManiaTopHandler,
		utils.FormatCommand("ctbtop"):   CTBCommandHandler,
		utils.FormatCommand("taikotop"): TaikoTopCommandHandler,
		utils.FormatCommand("map"):      MapCommandHandler,
		utils.FormatCommand("set-lol"):  SetLeagueCommandHandler,
		utils.FormatCommand("lol-curr"): CurrentLeagueGameCommand,
		utils.FormatCommand("rs-lol"):   RecentLeagueCommandHandler,
		utils.FormatCommand("lol"):      LeagueProfileCommandHandler,
		utils.FormatCommand("c"):        CompareCommandHandler,
		utils.FormatCommand("coc"):      CocProfileCommandHandler,
		utils.FormatCommand("coinflip"): CoinflipCommandHandler,
	}

	CommandsWOArgs = map[string]WOArgs{
		utils.FormatCommand("help"):          HelpCommandHandler,
		utils.FormatCommand("servers"):       ServersCommandHandler,
		utils.FormatCommand("lol-remove"):    RemoveLolCommandHandler,
		utils.FormatCommand("osu-remove"):    RemoveLolCommandHandler,
		utils.FormatCommand("balance"):       BalanceCommandHandler,
		utils.FormatCommand("reset-balance"): ResetBalance,
	}
}
