package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type EventRepo struct {
	ctx context.Context
	rdb *redis.Client
}

var (
	redisHost     string = utils.FindEnv("REDIS_HOST")
	redisPort     string = utils.FindEnv("REDIS_PORT")
	redisUsername string = utils.FindEnv("REDIS_USERNAME")
	redisPassword string = utils.FindEnv("REDIS_PASSWORD")
	redisDB       int    = 0
)

func NewEventRepo() *EventRepo {
	evr := &EventRepo{
		ctx: context.Background(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Username: redisUsername,
			Password: redisPassword,
			DB:       redisDB,
		}),
	}
	return evr
}

func (r *EventRepo) AddEndAuctionEvent(auction *models.Auction) error {
	expiration := time.Until(auction.EndTime)
	aid := auction.AuctionID.String()
	err := r.rdb.Set(r.ctx, "close_auction_ev::"+aid, aid, expiration).Err()
	if err != nil {
		return models.ErrCouldNotSetCloseEventForAuction
	}
	return nil
}

func (r *EventRepo) StartSubscriber(callback func(aid models.UUID), errCallback func(err error)) {
	pubsub := r.rdb.PSubscribe(r.ctx, "__keyevent@"+strconv.Itoa(redisDB)+"__:expired")
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(r.ctx)
		if err != nil {
			go errCallback(err)
			continue
		}

		if strings.HasPrefix(msg.Payload, "close_auction_ev::") {
			uuid := strings.TrimPrefix(msg.Payload, "close_auction_ev::")
			aid, err := utils.ParseUUID(uuid)
			if err != nil {
				go errCallback(err)
				continue
			}
			go callback(aid)
		}
	}
}
