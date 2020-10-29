package utils

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Champions struct {
	Type    string               `json:"type"`
	Format  string               `json:"format"`
	Version string               `json:"version"`
	Data    map[string]*Champion `json:"data"`
}

type Champion struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Blurb   string `json:"blurb"`
	Info    struct {
		Attack     int `json:"attack"`
		Defense    int `json:"defense"`
		Magic      int `json:"magic"`
		Difficulty int `json:"difficulty"`
	} `json:"info"`
	Image struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
	Tags    []string `json:"tags"`
	Partype string   `json:"partype"`
	Stats   struct {
		Hp                   float64 `json:"hp"`
		Hpperlevel           float64 `json:"hpperlevel"`
		Mp                   float64 `json:"mp"`
		Mpperlevel           float64 `json:"mpperlevel"`
		Movespeed            float64 `json:"movespeed"`
		Armor                float64 `json:"armor"`
		Armorperlevel        float64 `json:"armorperlevel"`
		Spellblock           float64 `json:"spellblock"`
		Spellblockperlevel   float64 `json:"spellblockperlevel"`
		Attackrange          float64 `json:"attackrange"`
		Hpregen              float64 `json:"hpregen"`
		Hpregenperlevel      float64 `json:"hpregenperlevel"`
		Mpregen              float64 `json:"mpregen"`
		Mpregenperlevel      float64 `json:"mpregenperlevel"`
		Crit                 float64 `json:"crit"`
		Critperlevel         float64 `json:"critperlevel"`
		Attackdamage         float64 `json:"attackdamage"`
		Attackdamageperlevel float64 `json:"attackdamageperlevel"`
		Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
		Attackspeed          float64 `json:"attackspeed"`
	} `json:"stats"`
}

type ChampionInfo struct {
	FreeChampionIDsForNewPlayers []int `json:"freeChampionIDsForNewPlayers"`
	FreeChampionIDs              []int `json:"freeChampionIDs"`
	MaxNewPlayerLevel            int   `json:"maxNewPlayerLevel"`
}

type ChampionMastery struct {
	ChestGranted                 bool   `json:"chestGranted"`
	ChampionLevel                int    `json:"championLevel"`
	ChampionPoints               int    `json:"championPoints"`
	ChampionID                   int    `json:"championId"`
	ChampionPointsUntilNextLevel int    `json:"championPointsUntilNextLevel"`
	LastPlayTime                 int    `json:"lastPlayTime"`
	TokensEarned                 int    `json:"tokensEarned"`
	ChampionPointsSinceLastLevel int    `json:"championPointsSinceLastLevel"`
	SummonerID                   string `json:"summonerId"`
}

type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountid"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileiconid"`
	Revisiondate  int64  `json:"revisiondate"`
	SummonerLevel int    `json:"summonerlevel"`
}

type SummonerRank struct {
	LeagueID     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
	MiniSeries   struct {
		Target   int    `json:"target"`
		Wins     int    `json:"wins"`
		Losses   int    `json:"losses"`
		Progress string `json:"progress"`
	} `json:"miniSeries,omitempty"`
}

type LiveMatch struct {
	GameID            int64         `json:"gameId"`
	MapID             int           `json:"mapId"`
	GameMode          string        `json:"gameMode"`
	GameType          string        `json:"gameType"`
	GameQueueConfigID int           `json:"gameQueueConfigId"`
	Participants      []Participant `json:"participants"`
	Observers         struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	PlatformID      string `json:"platformId"`
	BannedChampions []struct {
		ChampionID int `json:"championId"`
		TeamID     int `json:"teamId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bannedChampions"`
	GameStartTime int64 `json:"gameStartTime"`
	GameLength    int   `json:"gameLength"`
}

type Participant struct {
	TeamID                   uint8         `json:"teamId"`
	Spell1ID                 int           `json:"spell1Id"`
	Spell2ID                 int           `json:"spell2Id"`
	ChampionID               int           `json:"championId"`
	ProfileIconID            int           `json:"profileIconId"`
	SummonerName             string        `json:"summonerName"`
	Bot                      bool          `json:"bot"`
	SummonerID               string        `json:"summonerId"`
	GameCustomizationObjects []interface{} `json:"gameCustomizationObjects"`
	Perks                    struct {
		PerkIds      []int `json:"perkIds"`
		PerkStyle    int   `json:"perkStyle"`
		PerkSubStyle int   `json:"perkSubStyle"`
	} `json:"perks"`
}

type Matches struct {
	Matches []struct {
		PlatformID string `json:"platformId"`
		GameID     int64  `json:"gameId"`
		Champion   int    `json:"champion"`
		Queue      int    `json:"queue"`
		Season     int    `json:"season"`
		Timestamp  int64  `json:"timestamp"`
		Role       string `json:"role"`
		Lane       string `json:"lane"`
	} `json:"matches"`
	StartIndex int `json:"startIndex"`
	EndIndex   int `json:"endIndex"`
	TotalGames int `json:"totalGames"`
}

type SanitizedRank struct {
	SummonerName string
	Team         string
	Champion     string
	Solo         string
	Flex         string
}

type Matchlist struct {
	Matches    []*MatchReference `json:"matches"`
	TotalGames int               `json:"totalGames"`
	StartIndex int               `json:"startIndex"`
	EndIndex   int               `json:"endIndex"`
}

// MatchReference contains information about a game by a single summoner
type MatchReference struct {
	Lane       string `json:"lane"`
	GameID     int    `json:"gameId"`
	Champion   int    `json:"champion"`
	PlatformID string `json:"platformId"`
	Season     int    `json:"season"`
	Queue      int    `json:"queue"`
	Role       string `json:"role"`
	Timestamp  int    `json:"timestamp"`
}


type RiotClient struct {
	BaseURL   string `json:"base_url"`
	Token     string `json:"token"`
	Champions Champions
}

var ValidRegions map[string]string

func InitAPI() {
	ValidRegions = map[string]string{
		"euw": "euw1",
		"eun": "eun1",
		"br":  "br1",
		"kr":  "kr",
		"jp":  "jp1",
		"las": "la2",
		"lan": "la1",
		"na":  "na1",
		"oc":  "oc1",
		"tr":  "tr1",
		"ru":  "ru",
	}
}

// print out all the errors from a list
func LogErrors(errs []error) {
	for _, err := range errs {
		log.Print("err: ", err)
	}
}

// load champion with key, with this we can easily check summoner names
func (champions *Champions) GetChampionWithKey(key string) *Champion {
	for i := range champions.Data {
		if key == champions.Data[i].Key {
			return champions.Data[i]
		}
	}

	// return nil if the user is not found
	return nil
}

// Parse champion data from the riot api
func ParseChampions(sa *gorequest.SuperAgent) *Champions {
	champUrl := "http://ddragon.leagueoflegends.com/cdn/10.7.1/data/en_US/champion.json"
	var champs *Champions
	_, _, errs := sa.Clone().
		Get(champUrl).
		EndStruct(&champs)

	if errs != nil {
		LogErrors(errs)
	}

	return champs
}

// Check if a shorthand is a valid and return that actual string,
// also return a error, if there region is not valid
func CheckValidRegion(region string) (string, error) {
	value, ok := ValidRegions[strings.ToLower(region)]
	if !ok {
		return "", errors.New("region is not valid $lol-servers for all regions")
	}

	return value, nil
}

func NewRiotClient(region, token string) RiotClient {
	// we don't need to check if the key is valid, since that is done outside in another function
	sa := gorequest.New().Timeout(10*time.Second).
		Retry(2, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError)

	// load all the champions, so that we can use them later
	c := ParseChampions(sa)
	return RiotClient{
		BaseURL:   "https://" + region + ".api.riotgames.com/lol",
		Token:     token,
		Champions: *c,
	}
}

// NewAgent creates an returns a new request SuperAgent
func (c *RiotClient) NewAgent(path string, query string) *gorequest.SuperAgent {
	// form a request url using query parameters
	url := strings.Join([]string{c.BaseURL, path, query}, "/")

	// create the actual super agent using the riot token, and also set different retry parameters
	sa := gorequest.New().Get(url).Set("X-Riot-Token", c.Token).Timeout(10*time.Second).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		Retry(10, 5*time.Second, http.StatusTooManyRequests)

	return sa
}

func (c *RiotClient) GetSummonerWithName(name string) (*Summoner, error) {
	var summoner Summoner
	// create a new superAgent and request summoner information from api
	response, _, errs := c.NewAgent("summoner/v4/summoners/by-name", name).EndStruct(&summoner)
	if errs != nil {
		LogErrors(errs)
	}

	// check if the request was valid
	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return &summoner, nil
}

// Returns the solo&duo and flex ranks of a summoner with param: id
func (c *RiotClient) GetSummonerRankWithID(id string) ([]SummonerRank, error) {
	// array since summoners have 2 ranks: solo&duo and flex
	var rank []SummonerRank
	response, _, errs := c.NewAgent("league/v4/entries/by-summoner", id).EndStruct(&rank)
	if errs != nil {
		LogErrors(errs)
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return rank, nil
}

func (c *RiotClient) GetFreeRotation() (*ChampionInfo, error) {
	var freeRotation *ChampionInfo
	response, _, errs := c.NewAgent("platform/v3/champion-rotations", "").EndStruct(&freeRotation)
	if errs != nil {
		LogErrors(errs)
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return freeRotation, nil
}

func (c *RiotClient) GetSummonerLiveMatch(summoner *Summoner) (*LiveMatch, error) {
	var match LiveMatch
	response, _, errs := c.NewAgent("spectator/v4/active-games/by-summoner", summoner.ID).EndStruct(&match)

	if errs != nil {
		LogErrors(errs)
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return &match, nil
}

func (c *RiotClient) GetListOfMatches(accountId string, begin, end int) (*Matchlist, error) {
	var matches *Matchlist
	endpoint := fmt.Sprintf("matchlists/by-account/%s?beginIndex=%d&endIndex=%d",
		accountId, begin, end)
	response, _, errs := c.NewAgent(endpoint, "").EndStruct(&matches)
	if errs != nil {
		LogErrors(errs)
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return matches, nil
}

func (c *RiotClient) ListsSummonerMasteries(summonerID string) ([]*ChampionMastery, error) {
	var masteries []*ChampionMastery
	endpoint := fmt.Sprintf("champion-mastery/v4/champion-masteries/by-summoner/%s", summonerID)
	response, _, errs := c.NewAgent(endpoint, "").EndStruct(&masteries)
	if errs != nil {
		LogErrors(errs)
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return masteries, nil
}

func (c *RiotClient) GetSingleChampionMastery(summonerID, championID string) (
	*ChampionMastery, error) {
	var mastery *ChampionMastery
	endpoint := fmt.Sprintf("champion-mastery/v4/champion-masteries/by-summoner/%s/by-champion/%s",
		summonerID, championID)
	response, _, errs := c.NewAgent(endpoint, "").EndStruct(&mastery)
	if errs != nil {
		LogErrors(errs)
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return mastery, nil
}

// Sanitize user input into a table format
func (c *RiotClient) NewSanitizedRank(summonerName string, team uint8, championId int) SanitizedRank {
	t := "RED"
	if team == 100 {
		t = "BLUE"
	}

	// get the champion data from the preloaded champion data.
	champ := c.Champions.GetChampionWithKey(strconv.Itoa(championId))

	// Return N/A in the place of ranks, since this is filled later!
	return SanitizedRank{summonerName, t, champ.Name, "N/A", "N/A"}
}

func (c *RiotClient) GetLiveMatchBySummonerName(summonerName string) ([]SanitizedRank, error) {
	// find the summoner in question
	s, err := c.GetSummonerWithName(summonerName)
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	// get the live match, returns an error if the user is not in a match
	liveMatch, err := c.GetSummonerLiveMatch(s)
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	// create a wait group so that we can execute tasks concurrently, but still in an order
	wg := sync.WaitGroup{}
	wg.Add(len(liveMatch.Participants))
	participants := make([]SanitizedRank, len(liveMatch.Participants))
	for pi := range liveMatch.Participants {
		go func(i int, p Participant) {
			// make sure the wait counter is decreased b
			defer wg.Done()

			// get the summoner's rank
			r, err := c.GetSummonerRankWithID(p.SummonerID)
			if err != nil {
				log.Println("[ERROR]", err)
				return
			}

			// format the ranks
			sr := c.NewSanitizedRank(p.SummonerName, p.TeamID, p.ChampionID)
			for ri := range r {
				if r[ri].QueueType == "RANKED_SOLO_5x5" {
					sr.Solo = fmt.Sprintf("%s %s", r[ri].Tier, r[ri].Rank)
				}
				if r[ri].QueueType == "RANKED_FLEX_SR" {
					sr.Flex = fmt.Sprintf("%s %s", r[ri].Tier, r[ri].Rank)
				}
			}

			participants[i] = sr
		}(pi, liveMatch.Participants[pi])
	}

	// wait until the wait counter is 0, meaning all the processing is done
	wg.Wait()

	// sort the users so that they are grouped with their own teams
	sort.SliceStable(participants, func(a, b int) bool {
		return participants[a].Team < participants[b].Team
	})

	return participants, nil
}
