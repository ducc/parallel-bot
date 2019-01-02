package discord

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type Handlers struct {
	MessageCreateEvents chan *MessageCreateEvent
}

func NewHandlers() *Handlers {
	handlers := &Handlers{
		MessageCreateEvents: make(chan *MessageCreateEvent),
	}

	return handlers
}

func (h *Handlers) Register(discord *discordgo.Session) {
	discord.AddHandler(h.ReadyHandler)
	discord.AddHandler(h.MessageCreateHandler)
}

func logEvent(event interface{}) *log.Entry {
	return log.WithField("event", reflect.TypeOf(event).String())
}

func (h *Handlers) ReadyHandler(session *discordgo.Session, event *discordgo.Ready) {
	logEvent(event).Info()
}

func (h *Handlers) MessageCreateHandler(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Content == "" {
		return
	}

	logEvent(event).WithField("message", event.Message.Content).Debug()

	h.MessageCreateEvents <- &MessageCreateEvent{
		BaseEvent: BaseEvent{
			Session: session,
		},
		Event: event,
	}
}
