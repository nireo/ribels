package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	apiURLFormat      = "%s://%s.%s%s"
	baseURL           = "api.riotgames.com"
	scheme            = "https"
	apiTokenHeaderKey = "X-Riot-Token"
)

var client *Client

type Region string

type Requester struct {
	Response http.Response
	Respo
}

const (
	RegionBrasil            Region = "br1"
	RegionEuropeNorthEast          = "eun1"
	RegionEuropeWest               = "euw1"
	RegionJapan                    = "jp1"
	RegionKorea                    = "kr"
	RegionLatinAmericaNorth        = "la1"
	RegionLatinAmericaSouth        = "la2"
	RegionNorthAmerica             = "na1"
	RegionOceania                  = "oc1"
	RegionTurkey                   = "tr1"
	RegionRussia                   = "ru"
	RegionPBE                      = "pbe1"
)

var (
	Regions = []Region{
		RegionBrasil,
		RegionEuropeNorthEast,
		RegionEuropeWest,
		RegionJapan,
		RegionKorea,
		RegionLatinAmericaNorth,
		RegionLatinAmericaSouth,
		RegionNorthAmerica,
		RegionOceania,
		RegionTurkey,
		RegionRussia,
		RegionPBE,
	}
)

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
	return &request, nil
}
