package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Sharder struct {
	token       string
	ratelimiter *discordgo.RateLimiter
	shards      map[int]*Shard
	handlers    *Handlers
	mutex       *sync.Mutex
}

func NewSharder(token string, handlers *Handlers) *Sharder {
	return &Sharder{
		token:       token,
		ratelimiter: discordgo.NewRatelimiter(),
		shards:      make(map[int]*Shard, 1),
		handlers:    handlers,
		mutex:       &sync.Mutex{},
	}
}

func (s *Sharder) Shard(shardCount int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for shardID := 0; shardID < shardCount; shardID++ {
		shard, err := NewShard(s, shardCount, shardID)
		if err != nil {
			log.WithError(err).Error("error creating new shard")
			continue
		}

		if shard, ok := s.shards[shardID]; ok {
			if err := shard.Stop(); err != nil {
				log.WithError(err).Error("error stopping shard %d", shardID)
				continue
			}
		}

		s.shards[shardID] = shard

		if err := shard.Start(); err != nil {
			log.WithError(err).Error("error starting shard %d", shardID)
			continue
		}
	}

	for shardID, shard := range s.shards {
		if shardID >= shardCount {
			if err := shard.Stop(); err != nil {
				log.WithError(err).Error("error stopping shard %d", shardID)
				continue
			}
		}
	}
}

func (s *Sharder) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, shard := range s.shards {
		if err := shard.Stop(); err != nil {
			log.WithError(err).Error("error stopping shard %d", shard.session.ShardID)
			continue
		}
	}
}

type Shard struct {
	sharder *Sharder
	session *discordgo.Session
}

func NewShard(sharder *Sharder, shardCount, shardID int) (*Shard, error) {
	session, err := discordgo.New("Bot " + sharder.token)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating shard %d/%d", shardID, shardCount)
	}

	session.Ratelimiter = sharder.ratelimiter
	session.ShardCount, session.ShardID = shardCount, shardID

	sharder.handlers.Register(session)

	return &Shard{
		session: session,
	}, nil
}

func (s *Shard) Start() error {
	return s.session.Open()
}

func (s *Shard) Stop() error {
	return s.session.Close()
}
