package transforms

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type ExecutionResult struct {
	*CommandExecution
	result interface{}
}

func ExecuteCommands(commandExecutions <-chan *CommandExecution) <-chan *ExecutionResult {
	wg := &sync.WaitGroup{}
	executionResults := make(chan *ExecutionResult)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for commandExecution := range commandExecutions {
			result, err := commandExecution.command.Execute(commandExecution.session, commandExecution.event, commandExecution.args)
			if err != nil {
				log.WithError(err).WithField("name", commandExecution.command.Name()).Error("error executing command")
				continue
			}

			executionResults <- &ExecutionResult{
				CommandExecution: commandExecution,
				result:           result,
			}
		}
	}()

	go func() {
		wg.Wait()
		close(executionResults)
	}()

	return executionResults
}
