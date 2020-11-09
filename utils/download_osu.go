package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadOsuFile(beatmapID string) error {
	url := fmt.Sprintf("https://osu.ppy.sh/osu/%s", beatmapID)
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return err
	}

	file, err := os.Create(fmt.Sprintf("./temp/%s", beatmapID))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
