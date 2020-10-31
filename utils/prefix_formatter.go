package utils

var prefix string

func FormatCommand(command string) string {
	return prefix + command
}

func SetPrefix(p string) {
	prefix = p
}
