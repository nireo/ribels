package league_api

type Region string

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
