package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	oppai "github.com/flesnuk/oppai5"
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
	PPCountryRank string `json:"pp_country_rank"`
	RawPP         string `json:"pp_raw"`
	UserID        string `json:"user_id"`
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
	UserID      string `json:"user_id"`
}

type OsuBeatmap struct {
	Approved     string `json:"approved"`
	BPM          string `json:"bpm"`
	Difficulty   string `json:"difficultyrating"`
	Created      string `json:"creator"`
	Artist       string `json:"artist"`
	Title        string `json:"title"`
	TotalLength  string `json:"total_length"`
	MaxCombo     string `json:"max_combo"`
	Version      string `json:"version"`
	BeatmapSetID string `json:"beatmapset_id"`
}

type OsuRecentPlay struct {
	UserID      string `json:"user_id"`
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
	BeatmapID   string
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

type OsuMatch struct {
	MatchID   string `json:"match_id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type OsuGameScore struct {
	Slot        string `json:"slot"`
	Team        string `json:"team"`
	UserID      string `json:"user_id"`
	Score       string `json:"score"`
	MaxCombo    string `json:"max_combo"`
	Rank        string `json:"rank"`
	Count300    string `json:"count300"`
	Count100    string `json:"count100"`
	Count50     string `json:"count50"`
	CountMiss   string `json:"countmiss"`
	CountGeki   string `json:"countgeki"`
	CountKatu   string `json:"countkatu"`
	Pass        string `json:"pass"`
	EnabledMods string `json:"enabled_mods"`
}

type OsuGame struct {
	GameID      string `json:"game_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	BeatmapID   string `json:"beatmap_id"`
	PlayMode    string `json:"play_mode"`
	MatchType   string `json:"match_type"`
	ScoringType string `json:"scoring_type"`
	TeamType    string `json:"team_type"`
	Mods        string `json:"mods"`
	Scores      OsuGameScore
}

type OsuMultiMatch struct {
	Match OsuMatch
	Games []OsuGame
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

func GetMultiplayerGame(matchID string) (*OsuMultiMatch, error) {
	var multiMatch OsuMultiMatch
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_match?k=%s&mp=%s", key, matchID))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &multiMatch); err != nil {
		return nil, err
	}

	return &multiMatch, nil
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

func GetOsuBeatmapMods(beatmapID, mods string) (*OsuBeatmap, error) {
	// this object is used for returning an error without the risk of panicking
	singleMap := &OsuBeatmap{}
	var beatmaps []OsuBeatmap
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_beatmaps?b=%s&k=%s&limit=5&mods=%s", beatmapID, key, mods))
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

	requestURL := fmt.Sprintf("https://osu.ppy.sh/api/get_user_best?u=%s&k=%s&limit=3&m=%d",
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

	// predefine rank emojis
	RankEmojis = map[string]string{
		"X":  "<:bibelsX:753277439102418996>",
		"S":  "<:bibelsS:753277217420607679>",
		"XH": "<:bibelsXH:753277379048374334>",
		"SH": "<:bibelsSH:753277326128709665>",
		"A":  "<:bibelsA:753276834933637282>",
		"B":  "<:bibelsB:753276991473451020>",
		"C":  "<:bibelsC:753277059094020216>",
		"D":  "<:bibelsD:753277123070001244>",
		"F":  "<:bibelsF:776190735950544936>",
	}

	// create a mods string, so that we can easily find the mode number related to the mode name.
	mods = map[string]uint8{
		"standard": 0,
		"taiko":    1,
		"ctb":      2,
		"mania":    3,
	}
}

// Acc calculations for different plays, could be placed into different functions, but usage is more
// clear when used as independent methods
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

func (rp *OsuRecentPlay) CalculateAcc() string {
	missCount, _ := strconv.Atoi(rp.CountMiss)
	count50, _ := strconv.Atoi(rp.Count50)
	count100, _ := strconv.Atoi(rp.Count100)
	count300, _ := strconv.Atoi(rp.Count300)

	top := float64(50*count50 + 100*count100 + 300*count300)
	bot := float64(300 * (missCount + count300 + count100 + count50))
	acc := (top / bot) * 100

	return fmt.Sprintf("%.2f", acc)
}

func CalculateTaikoAcc(topPlay *OsuTopPlay) string {
	missCount, _ := strconv.ParseFloat(topPlay.CountMiss, 64)
	count100, _ := strconv.ParseFloat(topPlay.Count100, 64)
	count300, _ := strconv.ParseFloat(topPlay.Count300, 64)

	top := float64(0.5*count100 + count300)
	bot := float64(count300 + missCount + count100)
	acc := (top / bot) * 100

	return fmt.Sprintf("%.2f", acc)
}

func (rp *OsuScore) CalculateAcc() string {
	missCount, _ := strconv.Atoi(rp.CountMiss)
	count50, _ := strconv.Atoi(rp.Count50)
	count100, _ := strconv.Atoi(rp.Count100)
	count300, _ := strconv.Atoi(rp.Count300)

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

type MapResultPP struct {
	PlayPP float64 `json:"play_pp"`
	IfFCPP float64 `json:"if_fc"`
}

func (rp *OsuRecentPlay) CalculatePP() (*MapResultPP, error) {
	if err := DownloadOsuFile(rp.BeatmapID); err != nil {
		return &MapResultPP{}, errors.New("could not download .osu file")
	}

	file, err := os.Open(fmt.Sprintf("./temp/%s", rp.BeatmapID))
	if err != nil {
		return &MapResultPP{}, errors.New("could not parse file")
	}

	bmap := oppai.Parse(file)

	count300, _ := strconv.Atoi(rp.Count300)
	count100, _ := strconv.Atoi(rp.Count100)
	count50, _ := strconv.Atoi(rp.Count50)
	maxCombo, _ := strconv.Atoi(rp.MaxCombo)
	countMiss, _ := strconv.Atoi(rp.CountMiss)
	enabledMods, _ := strconv.Atoi(rp.EnabledMods)

	beatmap, err := GetOsuBeatmap(rp.BeatmapID)
	if err != nil {
		return &MapResultPP{}, err
	}
	maxMaxCombo, _ := strconv.Atoi(beatmap.MaxCombo)

	pp := oppai.PPInfo(bmap, &oppai.Parameters{
		N300:   uint16(count300),
		N100:   uint16(count100),
		N50:    uint16(count50),
		Misses: uint16(countMiss),
		Combo:  uint16(maxCombo),
		Mods:   uint32(enabledMods),
	}).PP

	ifFcpp := oppai.PPInfo(bmap, &oppai.Parameters{
		N300:   uint16(count300),
		N100:   uint16(count100 + countMiss),
		N50:    uint16(count50),
		Misses: 0,
		Combo:  uint16(maxMaxCombo),
		Mods:   uint32(enabledMods),
	}).PP

	// remove the file
	if err := os.Remove(fmt.Sprintf("./temp/%s", rp.BeatmapID)); err != nil {
		return &MapResultPP{}, err
	}

	result := &MapResultPP{
		IfFCPP: ifFcpp.Total,
		PlayPP: pp.Total,
	}

	return result, nil
}

func (rp *OsuScore) CalculatePP(currentMapID string) (*MapResultPP, error) {
	if err := DownloadOsuFile(currentMapID); err != nil {
		return &MapResultPP{}, errors.New("could not download .osu file")
	}

	file, err := os.Open(fmt.Sprintf("./temp/%s", currentMapID))
	if err != nil {
		return &MapResultPP{}, errors.New("could not parse file")
	}

	bmap := oppai.Parse(file)

	count300, _ := strconv.Atoi(rp.Count300)
	count100, _ := strconv.Atoi(rp.Count100)
	count50, _ := strconv.Atoi(rp.Count50)
	maxCombo, _ := strconv.Atoi(rp.MaxCombo)
	countMiss, _ := strconv.Atoi(rp.CountMiss)
	enabledMods, _ := strconv.Atoi(rp.EnabledMods)

	beatmap, err := GetOsuBeatmap(currentMapID)
	if err != nil {
		return &MapResultPP{}, err
	}
	maxMaxCombo, _ := strconv.Atoi(beatmap.MaxCombo)

	pp := oppai.PPInfo(bmap, &oppai.Parameters{
		N300:   uint16(count300),
		N100:   uint16(count100),
		N50:    uint16(count50),
		Misses: uint16(countMiss),
		Combo:  uint16(maxCombo),
		Mods:   uint32(enabledMods),
	}).PP

	ifFcpp := oppai.PPInfo(bmap, &oppai.Parameters{
		N300:   uint16(count300),
		N100:   uint16(count100 + countMiss),
		N50:    uint16(count50),
		Misses: 0,
		Combo:  uint16(maxMaxCombo),
		Mods:   uint32(enabledMods),
	}).PP

	// remove the file
	if err := os.Remove(fmt.Sprintf("./temp/%s", currentMapID)); err != nil {
		return &MapResultPP{}, err
	}

	result := &MapResultPP{
		IfFCPP: ifFcpp.Total,
		PlayPP: pp.Total,
	}

	return result, nil
}

func (topPlay *OsuTopPlay) CalculateDiff() (float64, error) {
	if err := DownloadOsuFile(topPlay.BeatmapID); err != nil {
		return 0, errors.New("could not download .osu file")
	}

	file, err := os.Open(fmt.Sprintf("./temp/%s", topPlay.BeatmapID))
	if err != nil {
		return 0, errors.New("could not parse file")
	}

	bmap := oppai.Parse(file)
	enabledMods, _ := strconv.Atoi(topPlay.EnabledMods)

	diff := oppai.PPInfo(bmap, &oppai.Parameters{
		Mods: uint32(enabledMods),
	}).Diff

	// remove the file
	if err := os.Remove(fmt.Sprintf("./temp/%s", topPlay.BeatmapID)); err != nil {
		return 0, err
	}

	return diff.Total, nil
}

// The same as 'CalculatePP' for osu scores, but this is optimized for checking many maps,
// so that we don't have to: request a map, make a file, do calculations, delete file. With this
// function we can just do calculations for each map
func PrecalculatePP(scores []OsuScore, currentMapID string) ([]*MapResultPP, error) {
	var mapResults []*MapResultPP

	// download the .osu file and save it in the ./temp directory, named with the currentMapID
	if err := DownloadOsuFile(currentMapID); err != nil {
		return mapResults, errors.New("could not download .osu file")
	}

	// load the beatmap, so that we can also check what the PP amount would be, for a full combo
	beatmap, err := GetOsuBeatmap(currentMapID)
	if err != nil {
		return mapResults, err
	}

	maxMaxCombo, _ := strconv.Atoi(beatmap.MaxCombo)

	// find the file, where the beatmap information is saved
	file, err := os.Open(fmt.Sprintf("./temp/%s", currentMapID))
	if err != nil {
		return mapResults, errors.New("could not parse file")
	}

	bmap := oppai.Parse(file)
	for _, score := range scores {
		count300, _ := strconv.Atoi(score.Count300)
		count100, _ := strconv.Atoi(score.Count100)
		count50, _ := strconv.Atoi(score.Count50)
		maxCombo, _ := strconv.Atoi(score.MaxCombo)
		countMiss, _ := strconv.Atoi(score.CountMiss)
		enabledMods, _ := strconv.Atoi(score.EnabledMods)

		pp := oppai.PPInfo(bmap, &oppai.Parameters{
			N300:   uint16(count300),
			N100:   uint16(count100),
			N50:    uint16(count50),
			Misses: uint16(countMiss),
			Combo:  uint16(maxCombo),
			Mods:   uint32(enabledMods),
		}).PP

		ifFcpp := oppai.PPInfo(bmap, &oppai.Parameters{
			N300:   uint16(count300),
			N100:   uint16(count100 + countMiss),
			N50:    uint16(count50),
			Misses: 0,
			Combo:  uint16(maxMaxCombo),
			Mods:   uint32(enabledMods),
		}).PP

		result := &MapResultPP{
			IfFCPP: ifFcpp.Total,
			PlayPP: pp.Total,
		}

		mapResults = append(mapResults, result)
	}

	// finally after doing the necessary calculations, we can delete the map,
	// since we don't want to save it
	if err := os.Remove(fmt.Sprintf("./temp/%s", currentMapID)); err != nil {
		return mapResults, err
	}

	return mapResults, nil
}
