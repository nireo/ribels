package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var cocKey string

// InitCocAPI sets the local cocKey used by the Clash of Clans api.
func InitCocAPI() {
	cocKey = os.Getenv("COC_API")
}

// COCPlayer is a golang struct that mimics the json response from the Clash of Clans api.
type COCPlayer struct {
	Clan struct {
		Tag       string `json:"tag"`
		ClanLevel int    `json:"clanLevel"`
		Name      string `json:"name"`
	} `json:"clan"`
	League struct {
		Name string `json:"name"`
	} `json:"league"`
	AttackWins       int    `json:"attackWins"`
	DefenseWins      int    `json:"defenseWins"`
	TownHallLevel    int    `json:"townHallLevel"`
	TownHallWeapon   int    `json:"townHallWeaponLevel"`
	WarStars         int    `json:"warStars"`
	Name             string `json:"name"`
	Tag              string `json:"tag"`
	BuilderHallLevel int    `json:"builderHallLevel"`
	Trophies         int    `json:"trophies"`
	Spells           []struct {
		Name     string `json:"name"`
		Level    int    `json:"level"`
		MaxLevel int    `json:"maxLevel"`
		Village  string `json:"village"`
	} `json:"spells"`
	Heroes []struct {
		Name     string `json:"name"`
		Level    int    `json:"level"`
		MaxLevel int    `json:"maxLevel"`
		Village  string `json:"village"`
	} `json:"heroes"`
	Troop []struct {
		Name     string `json:"name"`
		Level    int    `json:"level"`
		MaxLevel int    `json:"maxLevel"`
		Village  string `json:"village"`
	} `json:"troops"`
}

// GetCOCPlayerData returns the data of a user given a playerTag.
func GetCOCPlayerData(playerTag string) (*COCPlayer, error) {
	copy := playerTag[1:]
	copy = "%23" + copy

	var cocPlayer *COCPlayer
	endpoint := fmt.Sprintf("https://api.clashofclans.com/v1/players/%s", copy)
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + cocKey
	client := &http.Client{}
	request.Header.Set("authorization", bearer)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(strconv.Itoa(response.StatusCode))
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &cocPlayer); err != nil {
		return nil, err
	}

	return cocPlayer, nil
}
