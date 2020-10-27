package utils

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"errors"
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

type RiotClient struct {
	BaseURL   string `json:"base_url"`
	Token     string `json:"token"`
	Champions Champions
}

func (champions *Champions) GetChampionWithKey(key string) *Champion {
	for i := range champions.Data {
		if key == champions.Data[i].Key {
			return champions.Data[i]
		}
	}

	return nil
}

func ParseChampions(sa *gorequest.SuperAgent) *Champions {
	champUrl := "http://ddragon.leagueoflegends.com/cdn/10.7.1/data/en_US/champion.json"
	var champs *Champions
	_, _, errs := sa.Clone().
		Get(champUrl).
		EndStruct(&champs)

	if nil != errs {
		for i := range errs {
			log.Print("[ERROR] ", errs[i])
		}
	}

	return champs
}

func NewRiotClient(region string, token string, timeout int) RiotClient {
	sa := gorequest.New().Timeout(10*time.Second).
		Retry(2, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError)

	c := ParseChampions(sa)

	return RiotClient{
		BaseURL:   "https://" + region + ".api.riotgames.com/lol",
		Token:     token,
		Champions: *c,
	}
}

func (c *RiotClient) NewAgent(path string, query string) *gorequest.SuperAgent {
	url := strings.Join([]string{c.BaseURL, path, query}, "/")
	sa := gorequest.New().Get(url).Set("X-Riot-Token", c.Token).Timeout(10*time.Second).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		Retry(10, 5*time.Second, 429)
	return sa
}

func (c *RiotClient) GetSummonerWithName(name string) (*Summoner, error) {
	summoner := Summoner{}
	response, _, errs := c.NewAgent("summoner/v4/summoners/by-name", name).EndStruct(&summoner)

	if errs != nil {
		for _, err := range errs {
				log.Print("err: ", err)
		}
	}

	if 200 != response.StatusCode {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return &summoner, nil
}

func (c *RiotClient) GetSummonerRankWithID(id string) ([]SummonerRank, error) {
	var rank []SummonerRank
	response, _, errs := c.NewAgent("league/v4/entries/by-summoner", id).EndStruct(&rank)

	if errs != nil {
		for _, err := range errs {
			log.Println("err: ", err)
		}
	}

	if 200 != response.StatusCode {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return rank, nil
}

func (c *RiotClient) GetSummonerLiveMatch(summoner *Summoner) (*LiveMatch, error) {
	var match LiveMatch
	response, _, errs := c.NewAgent("spectator/v4/active-games/by-summoner", summoner.ID).EndStruct(&match)

	if errs != nil {
		for _, err := range errs {
			log.Println("err: ", err)
		}
	}

	if response.StatusCode != 200 {
		log.Println("err: ", response.Status)
		return nil, errors.New(response.Status)
	}

	return &match, nil
}

func (c *RiotClient) NewSanitizedRank(summonerName string, team uint8, championId int) SanitizedRank {
	t := "RED"
	if team == 100 {
		t = "BLUE"
	}
	champ := c.Champions.GetChampionWithKey(strconv.Itoa(championId))
	return SanitizedRank{summonerName, t, champ.Name, "N/A", "N/A"}
}

func (c *RiotClient) GetLiveMatchBySummonerName(summonerName *string) ([]SanitizedRank, error) {
	s, err := c.GetSummonerWithName(*summonerName)

	if err != nil {
		log.Println("err: ",  err)
		return nil, err
	}

	liveMatch, err := c.GetSummonerLiveMatch(s)
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(liveMatch.Participants))
	ps := make([]SanitizedRank, len(liveMatch.Participants))
	for pi := range liveMatch.Participants {
		go func(i int, p Participant) {
			defer wg.Done()
			r, err := c.GetSummonerRankWithID(p.SummonerID)

			if err != nil {
				log.Println("[ERROR]", err)
				return
			} else {
				sr := c.NewSanitizedRank(p.SummonerName, p.TeamID, p.ChampionID)
				for ri := range r {
					if r[ri].QueueType == "RANKED_SOLO_5x5" {
						sr.Solo = fmt.Sprintf("%s %s", r[ri].Tier, r[ri].Rank)
					}
					if r[ri].QueueType == "RANKED_FLEX_SR" {
						sr.Flex = fmt.Sprintf("%s %s", r[ri].Tier, r[ri].Rank)
					}
				}
				ps[i] = sr
			}
		}(pi, liveMatch.Participants[pi])
	}
	wg.Wait()
	sort.SliceStable(ps, func(i, j int) bool { return ps[i].Team < ps[j].Team })
	return ps, err
}