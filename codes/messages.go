package codes

var messages = map[string]string{
	"": "",
}

func getMessage(rawCode string) string {
	return messages[rawCode]
}
