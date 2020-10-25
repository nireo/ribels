package utils

import (
	"os"

	"github.com/KnutZuidema/golio"
)

var client *golio.Client

func GetClient() *golio.Client {
	return client
}

func InitClient() {
	client = golio.NewClient(os.Getenv("LEAGUE_API"))
}
