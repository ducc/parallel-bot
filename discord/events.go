package discord

import "github.com/bwmarrin/discordgo"

type BaseEvent struct {
	Session *discordgo.Session
}

type MessageCreateEvent struct {
	BaseEvent
	Event *discordgo.MessageCreate
}
