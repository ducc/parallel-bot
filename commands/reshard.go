package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ducc/parallel-bot/discord"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Reshard struct {
	Sharder *discord.Sharder
}

func (*Reshard) Name() string {
	return "!!!reshard"
}

func (r *Reshard) Execute(session *discordgo.Session, event *discordgo.MessageCreate, args []string) (interface{}, error) {
	logrus.Debug("reached reshard")

	if event.Author.Username != "spong" {
		return nil, nil
	}

	if len(args) == 0 {
		return "Usage: !!!reshard <num shards>", nil
	}

	shards, err := strconv.Atoi(args[0])
	if err != nil {
		return "Invalid argument!", nil
	}

	go func() {
		r.Sharder.Shard(shards)
	}()

	return nil, nil
}
