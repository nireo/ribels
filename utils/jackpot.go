package utils

import (
	"errors"
	"fmt"
)

type PlayerIdentity struct {
	DiscordID string `json:"discord_id"`
	Username  string `json:"username"`
}

var GameRunning = false
var Players map[PlayerIdentity]int64
var BetTotal int64 = 0

func AddPlayer(discordID, discordName string, wager int64) error {
	db := GetDatabase()
	var user EconomyUser

	// check that the user has a bank account
	if err := db.Where(&EconomyUser{DiscordID: discordID}).First(&user).Error; err != nil {
		return err
	}

	if wager < 1 {
		return errors.New("cannot bet negative values")
	}

	if user.Balance < wager {
		return errors.New("you don't have sufficient funds")
	}

	// decrease the amount of the wager, so that the user can't use money they don't have
	user.Balance -= wager
	db.Save(&user)
	BetTotal += wager

	playerIdentity := PlayerIdentity{
		Username:  discordName,
		DiscordID: discordID,
	}
	Players[playerIdentity] += wager

	return nil
}

func StartGame() {
	GameRunning = true
	Players = make(map[PlayerIdentity]int64)
}

func ClearGame() {
	for k := range Players {
		delete(Players, k)
	}

	GameRunning = false
	BetTotal = 0
}

func PrintPlayers() string {
	var players string
	for player, wager := range Players {
		// calculate the user's chanches of winning
		winChance := float64(wager) / float64(BetTotal) * 100

		players += fmt.Sprintf("`%s - %d (%.2f%%)`\n", player.Username, wager, winChance)
	}

	return players
}
