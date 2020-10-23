package league_api

import (
	"fmt"
	"log"
	"net/http"
)

const (
	apiURLFormat      = "%s://%s.%s%s"
	baseURL           = "api.riotgames.com"
	scheme            = "https"
	apiTokenHeaderKey = "X-Riot-Token"
)

var client *Client

type Client struct {
	Region Region `json:"region"`
	APIKey string `json:"api_key"`
}

func CreateNewClient(region Region, key string) *Client {
	return &Client{
		Region: region,
		APIKey: key,
	}
}

func (client *Client) ProcessRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	request, err := client.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

}

func (client *Client) NewRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method,
		fmt.Sprintf(apiURLFormat, scheme, client.Region, baseURL, endpoint), body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	request.Header.Add(apiTokenHeaderKey, client.APIKey)
	request.Header.Add("Accept", "application/json")
	return request, nil
}
