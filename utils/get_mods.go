package utils

import "strconv"

func GetMods(modNumber string) (string, error) {
	val, err := strconv.Atoi(modNumber)
	if err != nil {
		return "", err
	}
	mods := ""

	if val&(1<<0) != 0 {
		mods += "NF"
	}

	if val&(1<<1) != 0 {
		mods += "EZ"
	}

	if val&(1<<3) != 0 {
		mods += "HD"
	}

	if val&(1<<4) != 0 {
		mods += "HR"
	}

	if val&(1<<5) != 0 {
		mods += "SD"
	}

	if val&(1<<9) != 0 {
		mods += "NC"
	}

	if val&(1<<9) != 0 {
		mods += "NC"
	} else if val&(1<<6) != 0 {
		mods += "DT"
	}

	if val&(1<<7) != 0 {
		mods += "RX"
	}

	if val&(1<<8) != 0 {
		mods += "HT"
	}

	if val&(1<<10) != 0 {
		mods += "FL"
	}

	if val&(1<<12) != 0 {
		mods += "SO"
	}

	if val&(1<<14) != 0 {
		mods += "PF"
	}

	if mods == "" {
		mods = "No Mod"
	}

	return mods, nil
}
