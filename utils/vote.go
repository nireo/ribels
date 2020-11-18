package utils

type CurrentVote struct {
	Content   string
	Against   uint
	Agree     uint
	StartedBy string
}

var currentVote *CurrentVote

func GetCurrentVote() *CurrentVote {
	return currentVote
}

func StartVote(content, startedBy string) {
	currentVote = &CurrentVote{
		Content:   content,
		StartedBy: startedBy,
		Against:   0,
		Agree:     0,
	}
}
