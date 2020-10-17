package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nireo/ribels/utils"
)

func MessageHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// tokenize the input
	args := strings.Split(msg.Content, " ")
	if args[0] == "$set" {
		if len(args) == 1 {
			session.ChannelMessageSend(msg.ChannelID, "No username provided")
			return
		}

		osuName := utils.FormatName(args[1:])
		db := utils.GetDatabase()
		// check if name already in database
		var user utils.User
		if err := db.Where(&utils.User{OsuName: osuName}).First(&user).Error; err == nil {
			session.ChannelMessageSend(msg.ChannelID, "User already in database")
			return
		}

		// insert into database
		newUser := &utils.User{
			DiscordID: msg.Author.ID,
			OsuName:   osuName,
		}

		db.Create(&newUser)

		session.ChannelMessageSend(msg.ChannelID, "Saved user in database")
	}

	if args[0] == "$osu" {
		var osuName string
		if len(args) > 1 {
			osuName = utils.FormatName(args[1:])
		} else {
			user, err := utils.CheckIfSet(msg.Author.ID)
			if err != nil {
				session.ChannelMessageSend(msg.ChannelID, "Not set in database")
				return
			}

			osuName = user.OsuName
		}

		osuUserArray, err := utils.GetUserFromOSU(osuName)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		selectedUser := osuUserArray[0]

		// create embed fields
		var fields []*discordgo.MessageEmbedField
		fields = append(fields,
			&discordgo.MessageEmbedField{Name: "Playcount", Value: selectedUser.Playcount, Inline: false})

		fields = append(fields,
			&discordgo.MessageEmbedField{Name: "Rank", Value: selectedUser.PPRank, Inline: false})

		fields = append(fields,
			&discordgo.MessageEmbedField{Name: "Playtime", Value: selectedUser.SecondsPlayed, Inline: false})

		fields = append(fields,
			&discordgo.MessageEmbedField{Name: "Level", Value: selectedUser.Level, Inline: false})

		fields = append(fields,
			&discordgo.MessageEmbedField{Name: "Country", Value: selectedUser.Country, Inline: false})

		fields = append(fields,
			&discordgo.MessageEmbedField{Name: "Accuracy", Value: selectedUser.Accuracy, Inline: false})

		// create the final embed
		var messageEmbed discordgo.MessageEmbed
		messageEmbed.Title = utils.UnFormatName(selectedUser.Username)
		messageEmbed.Fields = fields
		messageEmbed.Type = "rich"

		session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
	}

	if args[0] == "$top" {
		var osuName string
		if len(args) > 1 {
			osuName = utils.FormatName(args[1:])
		} else {
			user, err := utils.CheckIfSet(msg.Author.ID)
			if err != nil {
				session.ChannelMessageSend(msg.ChannelID, "Not set in database")
				return
			}

			osuName = user.OsuName
		}

		topPlays, err := utils.GetUserTopplaysFromOSU(osuName)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		var fields []*discordgo.MessageEmbedField
		for index := range topPlays {
			beatmap, err := utils.GetOsuBeatmap(topPlays[index].BeatmapID)
			if err != nil {
				session.ChannelMessageSend(msg.ChannelID, err.Error())
				return
			}

			formattedPP := strings.Split(topPlays[index].PP, ".")
			formattedValue := fmt.Sprintf("PP: %s, Score set: %s",
				formattedPP[0], topPlays[index].Date)

			mods, err := utils.GetMods(topPlays[index].EnabledMods)
			if err != nil {
				session.ChannelMessageSend(msg.ChannelID, err.Error())
				return
			}

			formattedTitle := fmt.Sprintf("%s + %s", beatmap.Title, mods)
			fields = append(fields,
				&discordgo.MessageEmbedField{
					Name:   formattedTitle,
					Value:  formattedValue,
					Inline: false})
		}

		var messageEmbed discordgo.MessageEmbed
		messageEmbed.Title = fmt.Sprintf("osu! Standard top plays for %s", utils.UnFormatName(osuName))
		messageEmbed.Type = "rich"
		messageEmbed.Fields = fields

		session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
	}

	if args[0] == "$recent" {
		var osuName string
		if len(args) > 1 {
			osuName = utils.FormatName(args[1:])
		} else {
			user, err := utils.CheckIfSet(msg.Author.ID)
			if err != nil {
				session.ChannelMessageSend(msg.ChannelID, "Not set in database")
				return
			}

			osuName = user.OsuName
		}

		// get the most recent play from user
		recentPlay, err := utils.GetRecentPlay(osuName)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		var fields []*discordgo.MessageEmbedField
		formattedCounts := fmt.Sprintf("[%s/%s/%s/%s]",
			recentPlay.Count300, recentPlay.Count100, recentPlay.Count50, recentPlay.CountMiss)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Counts",
			Value:  formattedCounts,
			Inline: true,
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Score",
			Value:  recentPlay.Score,
			Inline: true,
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Date",
			Value:  recentPlay.Date,
			Inline: false,
		})

		beatmap, err := utils.GetOsuBeatmap(recentPlay.BeatmapID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		var messageEmbed discordgo.MessageEmbed
		messageEmbed.Title = beatmap.Title
		messageEmbed.Fields = fields

		session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("Most recent osu!Standard play for %s", utils.UnFormatName(osuName)))

		session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
	}
}
