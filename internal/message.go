package internal

type Message struct {
	IsSuccess bool
	Content   string
}

var messages = map[string]Message{
	"wiki-saved":      {true, "a wiki page saved successfully"},
	"wiki-removed":    {true, "a wiki page removed successfully"},
	"missing-comment": {false, "comment can't be empty"},
	"missing-path":    {false, "path can't be empty"},
	"save-error":      {false, "unexpected error while saving"},
}

func GetMessage(code string) Message {
	return messages[code]
}
