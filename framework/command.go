package framework

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string
	Execute(*discordgo.Session, *discordgo.MessageCreate, []string) (interface{}, error)
}
