package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type PlayerIdentity struct {
	DiscordID string `json:"discord_id"`
	Username  string `json:"username"`
	Wager     int64  `json:"wager"`
}

var GameRunning = false
var BetTotal int64 = 0
var PlayersArray []PlayerIdentity

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

	// check if the user is already in the game
	index, found := InPlayers(discordID)
	if !found {
		playerIdentity := PlayerIdentity{
			Username:  discordName,
			DiscordID: discordID,
			Wager:     wager,
		}

		PlayersArray = append(PlayersArray, playerIdentity)
	} else {
		PlayersArray[index].Wager += wager
	}

	BetTotal += wager
	return nil
}

// Return the index of the player, and a boolean value if they were found in the array
func InPlayers(discordID string) (int, bool) {
	for index, player := range PlayersArray {
		if player.DiscordID == discordID {
			return index, true
		}
	}

	return 0, false
}

type PlayerChance struct {
	Player    PlayerIdentity
	TopTicket float64
	MinTicket float64
}

// Return the Player struct with discord id and username, and also returns the winning ticket
func ChooseWinner() (PlayerIdentity, float64) {
	var playersWithChances []PlayerChance
	playersWithChances = append(playersWithChances, PlayerChance{
		MinTicket: 0.0,
		TopTicket: float64(PlayersArray[0].Wager) / float64(BetTotal),
		Player:    PlayersArray[0],
	})

	// construct the ticket treshold array (skip the first element, since already in array)
	for index, player := range PlayersArray[1:] {
		playersWithChances = append(playersWithChances, PlayerChance{
			MinTicket: playersWithChances[index].TopTicket + 0.0000001,
			TopTicket: playersWithChances[index].TopTicket + (float64(player.Wager) / float64(BetTotal)),
			Player:    PlayersArray[index],
		})
	}

	// create the winning ticket
	rand.Seed(time.Now().UnixNano())
	winning := rand.Float64()

	for _, player := range playersWithChances {
		if player.MinTicket <= winning && winning <= player.TopTicket {
			// pay the winner
			var user EconomyUser
			db.Where(&EconomyUser{DiscordID: player.Player.DiscordID}).First(&user)
			user.Balance += BetTotal
			db.Save(&user)

			return player.Player, winning
		}
	}

	return PlayerIdentity{}, winning
}

func StartGame() {
	GameRunning = true
}

func ClearGame() {
	PlayersArray = nil
	GameRunning = false
	BetTotal = 0
}

func PrintPlayers() string {
	var players string
	players += fmt.Sprintf("Current jackpot is: `%d`\n", BetTotal)
	for _, player := range PlayersArray {
		// calculate the user's chanches of winning
		winChance := float64(player.Wager) / float64(BetTotal) * 100

		players += fmt.Sprintf("`%s - %d (%.2f%%)`\n", player.Username, player.Wager, winChance)
	}

	return players
}
