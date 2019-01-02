package commands

import "github.com/bwmarrin/discordgo"

type Ping struct{}

func (*Ping) Name() string {
	return "!!!ping"
}

func (*Ping) Execute(session *discordgo.Session, event *discordgo.MessageCreate, args []string) (interface{}, error) {
	return "Pong!", nil
}
