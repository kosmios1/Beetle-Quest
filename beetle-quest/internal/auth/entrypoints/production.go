//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"
	"crypto/tls"
	"strconv"

	"beetle-quest/pkg/models"
	arepo "beetle-quest/pkg/repositories/impl/http/admin"
	urepo "beetle-quest/pkg/repositories/impl/http/user"
	"beetle-quest/pkg/utils"

	"github.com/go-session/redis/v3"
	"github.com/go-session/session/v3"
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

func NewAuthController() (*controller.AuthController, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, models.ErrCouldNotLoadClientCetrificate
	}
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			TLSConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				ServerName:         redisHost,
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			},
			Addr:     redisHost + ":" + redisPort,
			Password: redisPassword,
			DB:       redisDB,
		})),
	)
	return controller.NewAuthController(
		service.NewAuthService(
			urepo.NewUserRepo(),
			arepo.NewAdminRepo()),
	), nil
}
