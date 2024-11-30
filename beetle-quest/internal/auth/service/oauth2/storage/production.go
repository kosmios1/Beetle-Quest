//go:build !beetleQuestTest

package storage

import (
	"beetle-quest/pkg/utils"
	"strconv"

	"github.com/go-oauth2/oauth2/v4"

	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
)

var (
	redisHost     string = utils.FindEnv("REDIS_HOST")
	redisPort     string = utils.FindEnv("REDIS_PORT")
	redisPassword string = utils.FindEnv("REDIS_PASSWORD_OAUTH2")
	redisUsername string = utils.FindEnv("REDIS_USERNAME_OAUTH2")
	redisDB       int    = utils.PanicIfError[int](strconv.Atoi(utils.FindEnv("REDIS_DB_OAUTH2")))
)

func GetTokenStorage() oauth2.TokenStore {
	return oredis.NewRedisStore(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Username: redisUsername,
		Password: redisPassword,
		DB:       redisDB,
	})
}
