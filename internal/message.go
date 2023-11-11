package internal

type Message struct {
	IsSuccess bool
	Content   string
}

var messages = map[string]Message{
	"ws": {true, "a wiki page saved successfully"},   // Wiki Saved
	"wd": {true, "a wiki page removed successfully"}, // Wiki deleted
	"mc": {false, "comment can't be empty"},          // Missing Comment
}

func GetMessage(code string) Message {
	return messages[code]
}
