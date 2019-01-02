package transforms

import (
	log "github.com/sirupsen/logrus"
)

func HandleExecutionResults(executionResults <-chan *ExecutionResult) {
	go func() {
		for commandExecution := range executionResults {
			if commandExecution.result == nil {
				continue
			}

			switch commandExecution.result.(type) {
			case string:
				if _, err := commandExecution.session.ChannelMessageSend(commandExecution.event.ChannelID, commandExecution.result.(string)); err != nil {
					log.WithError(err).WithField("message", commandExecution.result).WithField("channel", commandExecution.event.ChannelID).Error("error sending message to channel")
				}
			}
		}
	}()
}
