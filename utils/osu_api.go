package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var key string

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
	BeatmapID       string `json:"beatmap_id"`
	ScoreID         string `json:"score_id"`
	Score           string `json:"score"`
	MaxCombo        string `json:"maxcombo"`
	Count50         string `json:"count50"`
	Count100        string `json:"count100"`
	Count300        string `json:"count300"`
	CountMiss       string `json:"countmiss"`
	CountKatu       string `json:"countkatu"`
	CountGeki       string `json:"countgeki"`
	Perfect         string `json:"perfect"`
	EnabledMods     string `json:"enabled_mods"`
	Date            string `json:"date"`
	Rank            string `json:"rank"`
	PP              string `json:"pp"`
	ReplayAvailable string `json:"replay_available"`
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
}

func GetUserFromOSU(username string) ([]OsuUserResponse, error) {
	var osuUser []OsuUserResponse
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_user?u=%s&k=%s", username, key))
	if err != nil {
		return osuUser, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return osuUser, err
	}

	if err := json.Unmarshal(body, &osuUser); err != nil {
		return osuUser, err
	}

	return osuUser, nil
}

func GetUserTopplaysFromOSU(username string) ([]OsuTopPlay, error) {
	var topplays []OsuTopPlay
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_user_best?u=%s&k=%s", username, key))
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

func GetOsuBeatmap(beatmapID string) (OsuBeatmap, error) {
	var beatmap []OsuBeatmap
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_beatmaps?b=%s&k=%s&limit=5", beatmapID, key))
	if err != nil {
		return beatmap[0], err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return beatmap[0], err
	}

	if err := json.Unmarshal(body, &beatmap); err != nil {
		return beatmap[0], err
	}

	return beatmap[0], nil
}

func GetRecentPlay(username string) (OsuRecentPlay, error) {
	var beatmap []OsuRecentPlay
	response, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/api/get_beatmaps?u=%s&k=%s&limit=1", username, key))
	if err != nil {
		return beatmap[0], err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return beatmap[0], err
	}

	if err := json.Unmarshal(body, &beatmap); err != nil {
		return beatmap[0], err
	}

	return beatmap[0], nil
}

func InitApiKey() {
	key = os.Getenv("OSU_KEY")
}
