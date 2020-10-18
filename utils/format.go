package utils

import "strings"

// Result: ["install", "gentoo"] => "install_gentoo" | This the osu api format
func FormatName(name []string) string {
	return strings.Join(name, "_")
}

// Result: "instalL_gentoo" => "install gentoo" | Just to make the display prettier
func UnFormatName(formattedName string) string {
	return strings.Join(strings.Split(formattedName, "_"), " ")
}
