package utils

import (
	"github.com/KnutZuidema/golio/api"
)

var Servers map[string]api.Region

func InitClient() {
	Servers = map[string]api.Region{
		"euw": api.RegionEuropeWest,
		"eun": api.RegionEuropeNorthEast,
		"br":  api.RegionBrasil,
		"jp":  api.RegionJapan,
		"kr":  api.RegionKorea,
		"lan": api.RegionLatinAmericaNorth,
		"las": api.RegionLatinAmericaSouth,
		"oc":  api.RegionOceania,
		"ru":  api.RegionRussia,
		"pbe": api.RegionPBE,
		"tr":  api.RegionTurkey,
		"na":  api.RegionNorthAmerica,
	}
}
