package utils

import "strings"

func FormatName(name []string) string {
	return strings.Join(name, "_")
}

func UnFormatName(formattedName string) string {
	return strings.Join(strings.Split(formattedName, "_"), " ")
}
