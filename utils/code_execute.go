package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type CodeExecutionResponse struct {
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	Output     string `json:"output"`
	CPUTime    string `json:"cpuTime"`
}

func ExecuteCodeRequest(lang, code string) (CodeExecutionResponse, error) {
	var executionInfo CodeExecutionResponse
	reqBody, err := json.Marshal(map[string]string{
		"clientId":     os.Getenv("CODE_CLIENT_ID"),
		"clientSecret": os.Getenv("CODE_CLIENT_SECRET"),
		"language":     lang,
		"script":       code,
		"versionIndex": "3",
	})
	if err != nil {
		return executionInfo, err
	}

	response, err := http.Post("https://api.jdoodle.com/v1/execute", "application/json",
		bytes.NewBuffer(reqBody))

	if err != nil {
		return executionInfo, err
	}

	if response.StatusCode != 200 {
		return executionInfo, errors.New(strconv.Itoa(response.StatusCode))
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return executionInfo, err
	}

	if err := json.Unmarshal(body, &executionInfo); err != nil {
		return executionInfo, err
	}

	return executionInfo, nil
}
