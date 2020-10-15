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

		osuName := strings.Join(args[1:], "_")
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
		user, err := utils.CheckIfSet(msg.Author.ID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Not set in database")
			return
		}

		osuUserArray, err := utils.GetUserFromOSU(user.OsuName)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		selectedUser := osuUserArray[0]

		// create embed fields
		var fields []*discordgo.MessageEmbedField
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Playcount", Value: selectedUser.Playcount, Inline: false})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Rank", Value: selectedUser.PPRank, Inline: false})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Playtime", Value: selectedUser.SecondsPlayed, Inline: false})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Level", Value: selectedUser.Level, Inline: false})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Country", Value: selectedUser.Country, Inline: false})
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Accuracy", Value: selectedUser.Accuracy, Inline: false})

		// create the final embed
		var messageEmbed discordgo.MessageEmbed
		messageEmbed.Title = selectedUser.Username
		messageEmbed.Fields = fields
		messageEmbed.Type = "rich"

		session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
	}

	if args[0] == "$top" {
		user, err := utils.CheckIfSet(msg.Author.ID)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Not set in database")
			return
		}

		topPlays, err := utils.GetUserTopplaysFromOSU(user.OsuName)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, err.Error())
			return
		}

		var fields []*discordgo.MessageEmbedField
		for index := range topPlays {
			fields = append(fields,
				&discordgo.MessageEmbedField{
					Name:  "YEP",
					Value: fmt.Sprintf("BeatmapId %s and pp %s", topPlays[index].BeatmapID, topPlays[index].PP)})
		}

		var messageEmbed discordgo.MessageEmbed
		messageEmbed.Title = fmt.Sprintf("%s top plays", user.OsuName)
		messageEmbed.Type = "rich"
		messageEmbed.Fields = fields

		session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
	}
}
