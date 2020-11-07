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

var key string
var currentMap string

// Since golang is a static typed language we need to create structs for the json requests
// note that we don't need to add a field for every single json field, just for those which we need
type OsuUserResponse struct {
	Username      string `json:"username"`
	Playcount     string `json:"playcount"`
	RankedScore   string `json:"ranked_score"`
	PPRank        string `json:"pp_rank"`
	Level         string `json:"level"`
	Accuracy      string `json:"accuracy"`
	Country       string `json:"country"`
	SecondsPlayed string `json:"total_seconds_played"`
}

type OsuTopPlay struct {
	BeatmapID   string `json:"beatmap_id"`
	ScoreID     string `json:"score_id"`
	Score       string `json:"score"`
	MaxCombo    string `json:"maxcombo"`
	Count50     string `json:"count50"`
	Count100    string `json:"count100"`
	Count300    string `json:"count300"`
	CountMiss   string `json:"countmiss"`
	Perfect     string `json:"perfect"`
	EnabledMods string `json:"enabled_mods"`
	Date        string `json:"date"`
	Rank        string `json:"rank"`
	PP          string `json:"pp"`
}

type OsuBeatmap struct {
	Approved    string `json:"approved"`
	BPM         string `json:"bpm"`
	Difficulty  string `json:"difficultyrating"`
	Created     string `json:"creator"`
	Artist      string `json:"artist"`
	Title       string `json:"title"`
	TotalLength string `json:"total_length"`
	MaxCombo    string `json:"max_combo"`
	Version     string `json:"version"`
}

type OsuRecentPlay struct {
	BeatmapID   string `json:"beatmap_id"`
	Score       string `json:"score"`
	MaxCombo    string `json:"maxcombo"`
	Count50     string `json:"count50"`
	Count100    string `json:"count100"`
	Count300    string `json:"count300"`
	CountMiss   string `json:"countmiss"`
	EnabledMods string `json:"enabled_mods"`
	Date        string `json:"date"`
	Rank        string `json:"rank"`
}

type OsuScore struct {
	ScoreID     string `json:"score_id"`
	Score       string `json:"score"`
	Username    string `json:"username"`
	Count300    string `json:"count300"`
	Count100    string `json:"count100"`
	Count50     string `json:"count50"`
	CountMiss   string `json:"countmiss"`
	EnabledMods string `json:"enabled_mods"`
	UserID      string `json:"user_id"`
	Date        string `json:"date"`
	Rank        string `json:"rank"`
	PP          string `json:"pp"`
	MaxCombo    string `json:"maxcombo"`
}

func SetCurrentMap(currMapID string) {
	currentMap = currMapID
}

func GetCurrentMap() string {
	return currentMap
}

var RankEmojis map[string]string

func GetUserFromOSU(username string) (*OsuUserResponse, error) {
	// The osu api returns an array for some reason
	osuUser := &OsuUserResponse{}
	var osuUsers []OsuUserResponse
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_user?u=%s&k=%s", username, key))
	if err != nil {
		return osuUser, err
	}

	defer response.Body.Close()

	// read the body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return osuUser, err
	}

	// parse the json to golang structs
	if err := json.Unmarshal(body, &osuUsers); err != nil {
		return osuUser, err
	}

	osuUser = &osuUsers[0]
	return osuUser, nil
}

func GetPlaysInBeatmapFromUser(beatmapID, userID string) ([]OsuScore, error) {
	var osuScores []OsuScore
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_scores?u=%s&k=%s&limit=3&b=%s",
		userID, key, beatmapID))
	if err != nil {
		return osuScores, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return osuScores, err
	}

	if err := json.Unmarshal(body, &osuScores); err != nil {
		return osuScores, err
	}

	return osuScores, nil
}

func GetScoresForCurrentMap(username string) ([]OsuScore, error) {
	if currentMap == "" {
		return nil, errors.New("no current map")
	}

	var osuScores []OsuScore
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_scores?u=%s&k=%s&limit=3&b=%s",
		username, key, currentMap))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return osuScores, err
	}

	if err := json.Unmarshal(body, &osuScores); err != nil {
		return nil, err
	}

	return osuScores, nil
}

func GetOsuBeatmap(beatmapID string) (*OsuBeatmap, error) {
	// this object is used for returning an error without the risk of panicking
	singleMap := &OsuBeatmap{}
	var beatmaps []OsuBeatmap
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_beatmaps?b=%s&k=%s&limit=5", beatmapID, key))
	if err != nil {
		return singleMap, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return singleMap, err
	}

	if err := json.Unmarshal(body, &beatmaps); err != nil {
		return singleMap, err
	}

	singleMap = &beatmaps[0]

	return singleMap, nil
}

var mods map[string]uint8

func GetModeTopPlays(username, mode string) ([]OsuTopPlay, error) {
	var topplays []OsuTopPlay

	requestURL := fmt.Sprintf("https://osu.ppy.sh/api/get_user_best?u=%s&k=%s&limit=5&m=%d",
		username, key, mods[mode])

	response, err := http.Get(requestURL)
	if err != nil {
		return topplays, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return topplays, err
	}
	if err := json.Unmarshal(body, &topplays); err != nil {
		return topplays, err
	}

	return topplays, nil
}

func GetRecentPlay(username string) (*OsuRecentPlay, error) {
	singleMap := &OsuRecentPlay{}
	var beatmaps []OsuRecentPlay
	response, err := http.Get(fmt.Sprintf(
		"https://osu.ppy.sh/api/get_user_recent?u=%s&k=%s", username, key))

	if err != nil {
		return singleMap, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return singleMap, err
	}

	if err := json.Unmarshal(body, &beatmaps); err != nil {
		return singleMap, err
	}

	if len(beatmaps) == 0 {
		return singleMap, errors.New("user has no recent plays")
	}

	singleMap = &beatmaps[0]
	return singleMap, nil
}

func InitApiKey() {
	key = os.Getenv("OSU_KEY")

	RankEmojis = map[string]string{
		"X":  "<:bibelsX:753277439102418996>",
		"S":  "<:bibelsS:753277217420607679>",
		"XH": "<:bibelsXH:753277379048374334>",
		"SH": "<:bibelsSH:753277326128709665>",
		"A":  "<:bibelsA:753276834933637282>",
		"B":  "<:bibelsB:753276991473451020>",
		"C":  "<:bibelsC:753277059094020216>",
		"D":  "<:bibelsD:753277123070001244>",
		"F":  "kantsii lisää F emote ;)",
	}

	mods = map[string]uint8{
		"standard": 0,
		"taiko":    1,
		"ctb":      2,
		"mania":    3,
	}
}

func (tp *OsuTopPlay) CalculateTopPlayAcc() string {
	// format all the counts into numbers
	missCount, _ := strconv.Atoi(tp.CountMiss)
	count50, _ := strconv.Atoi(tp.Count50)
	count100, _ := strconv.Atoi(tp.Count100)
	count300, _ := strconv.Atoi(tp.Count300)

	top := float64(50*count50 + 100*count100 + 300*count300)
	bot := float64(300 * (missCount + count300 + count100 + count50))
	acc := (top / bot) * 100

	return fmt.Sprintf("%.2f", acc)
}

func GetOsuUsername(discordId string, args []string) (string, error) {
	var osuName string
	if len(args) > 1 {
		osuName = FormatName(args[1:])
	} else {
		user, err := CheckIfSet(discordId)
		if err != nil {
			return osuName, errors.New("could not find user")
		}

		osuName = user.OsuName
	}

	return osuName, nil
}
