package utils

import "strings"

// FormatName formats the given name array to ["install", "gentoo"] => "install_gentoo" | This the osu api format
func FormatName(name []string) string {
	return strings.Join(name, "_")
}

// UnFormatName is the inverse of the FormatName function.
func UnFormatName(formattedName string) string {
	return strings.Join(strings.Split(formattedName, "_"), " ")
}
