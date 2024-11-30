package repository

import (
  "beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"context"
	"errors"
	"strconv"
	"time"

	"crypto/tls"
	"github.com/go-redis/redis/v8"
)

var (
	redisHost     string = utils.FindEnv("REDIS_HOST")
	redisPort     string = utils.FindEnv("REDIS_PORT")
	redisPassword string = utils.FindEnv("REDIS_PASSWORD")
	redisUsername string = utils.FindEnv("REDIS_USERNAME")
	redisDB       int    = utils.PanicIfError[int](strconv.Atoi(utils.FindEnv("REDIS_DB_SESSION")))
  certFile             = "/serverCert.pem"
  keyFile              = "/serverKey.pem"
)

type SessionRepo struct {
	client *redis.Client
}

func NewSessionRepo() (*SessionRepo, error) {
  cert, err := tls.LoadX509KeyPair(certFile, keyFile)
  if err != nil {
      return nil, models.ErrCouldNotLoadClientCetrificate
  }
	client := redis.NewClient(&redis.Options{
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
	return &SessionRepo{client: client}, nil
}

func (s *SessionRepo) CreateSession(token string) error {
	ctx := context.Background()
	now := time.Now()
	end := now.Add(time.Hour * 24)
	err := s.client.Set(ctx, token, token, end.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionRepo) RevokeToken(token string) error {
	ctx := context.Background()
	err := s.client.Del(ctx, token).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionRepo) FindToken(token string) (string, error) {
	ctx := context.Background()
	val, err := s.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return "", errors.New("token not found")
	} else if err != nil {
		return "", err
	}
	return val, nil
}
