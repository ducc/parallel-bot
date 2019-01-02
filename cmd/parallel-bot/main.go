package main

import (
	"flag"
	"github.com/ducc/parallel-bot/commands"
	"github.com/ducc/parallel-bot/discord"
	"github.com/ducc/parallel-bot/framework"
	"github.com/ducc/parallel-bot/transforms"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	logLevel string
	token    string
)

func init() {
	flag.StringVar(&logLevel, "log-level", "info", "logrus logging level")
	flag.StringVar(&token, "token", "", "discord bot token")
}

func main() {
	flag.Parse()

	if ll, err := log.ParseLevel(logLevel); err != nil {
		log.WithError(err).WithField("logLevel", logLevel).Fatal("error parsing logrus logging level")
	} else {
		log.SetLevel(ll)
	}

	handlers := discord.NewHandlers()
	sharder := discord.NewSharder(token, handlers)

	registry := framework.NewRegistry()
	registry.Register(&commands.Ping{})
	registry.Register(&commands.Reshard{sharder})

	commandExecutions := transforms.ToCommandExecutions(handlers.MessageCreateEvents, registry)
	executionResults := transforms.ExecuteCommands(commandExecutions)
	transforms.HandleExecutionResults(executionResults)

	sharder.Shard(1)
	defer sharder.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigChan
}
