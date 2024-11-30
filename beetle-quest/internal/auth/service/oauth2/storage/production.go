//go:build !beetleQuestTest

package storage

import (
	"beetle-quest/pkg/utils"
	"strconv"
	"crypto/tls"

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
  certFile             = "/serverCert.pem"
  keyFile              = "/serverKey.pem"
)

func GetTokenStorage() oauth2.TokenStore {
  cert := utils.PanicIfError[tls.Certificate](tls.LoadX509KeyPair(certFile, keyFile))

	return oredis.NewRedisStore(&redis.Options{
	  TLSConfig: &tls.Config{
      MinVersion: tls.VersionTLS12,
      ServerName: redisHost,
      InsecureSkipVerify: true,
      Certificates: []tls.Certificate{cert},
	  },
		Addr:     redisHost + ":" + redisPort,
		Username: redisUsername,
		Password: redisPassword,
		DB:       redisDB,
	})
}
