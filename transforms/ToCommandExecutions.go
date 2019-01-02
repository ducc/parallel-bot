package transforms

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ducc/parallel-bot/discord"
	"github.com/ducc/parallel-bot/framework"
	"strings"
	"sync"
)

type CommandExecution struct {
	command framework.Command
	session *discordgo.Session
	event   *discordgo.MessageCreate
	args    []string
}

func ToCommandExecutions(events <-chan *discord.MessageCreateEvent, registry *framework.Registry) <-chan *CommandExecution {
	wg := &sync.WaitGroup{}
	commandExecutions := make(chan *CommandExecution)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for event := range events {
			if event.Event.Content == "" {
				continue
			}

			splits := strings.Split(event.Event.Content, " ")
			if len(splits) == 0 {
				continue
			}

			var args []string
			if len(splits) > 0 {
				args = splits[1:]
			}

			command := registry.Get(splits[0])
			if command == nil {
				continue
			}

			commandExecutions <- &CommandExecution{
				command: command,
				session: event.Session,
				event:   event.Event,
				args:    args,
			}
		}
	}()

	go func() {
		wg.Wait()
		close(commandExecutions)
	}()

	return commandExecutions
}
