package application

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	pkgClickHouse "hezzl/_pkg/clickhouse"
	"hezzl/_pkg/natsPkg"
	"hezzl/_pkg/postgresql"
	"hezzl/internal/application/config"
	"hezzl/internal/application/product"
	"hezzl/internal/domain/event"
	"os"
	"strings"
)

type App interface {
	Run(ctx context.Context) error
}

type app struct {
	logEntry    logrus.Entry
	controllers *product.Controllers
}

type Config struct {
	ProjectAbsPath string `env:"PROJECT_ABS_PATH" env-required:"true"`

	PostgresHost     string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort     string `env:"POSTGRES_PORT" env-required:"true"`
	PostgresUsername string `env:"POSTGRES_USERNAME" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDBName   string `env:"POSTGRES_DB_NAME" env-required:"true"`

	RedisAddr              string `env:"REDIS_ADDR" env-required:"true"`
	RedisPassword          string `env:"REDIS_PASSWORD" env-required:"true"`
	RedisDB                int    `env:"REDIS_DB" env-required:"true"`
	RedisExpirationSeconds int    `env:"REDIS_EXPIRATION_SECONDS" env-required:"true"`

	APIHost string `env:"API_HOST" env-required:"true"`
	APIPort string `env:"API_PORT" env-required:"true"`

	NatsURL        string `env:"NATS_URL" env-required:"true"`
	NatsFetchTime  int    `env:"NATS_FETCH_TIME" env-required:"true"`
	GoodsChangeKey string `env:"GOODS_CHANGE_KEY" env-required:"true"`

	ClickhouseAddr     string `env:"CLICKHOUSE_ADDR" env-required:"true"`
	ClickhouseUsername string `env:"CLICKHOUSE_USERNAME" env-required:"true"`
	ClickhousePassword string `env:"CLICKHOUSE_PASSWORD" env-required:"true"`
	ClickhouseDBName   string `env:"CLICKHOUSE_DB_NAME" env-required:"true"`
}

func ReadEnv(logger *logrus.Entry) (*Config, error) {
	c := Config{}
	envFilePath := os.Getenv("ENV_FILE_PATH")
	if strings.TrimSpace(envFilePath) == "" {
		logger.Warnln("env variable ENV_FILE_PATH is not provided so ENV_FILE_PATH set to default \"./.env\"")
		envFilePath = "./.env"
	}
	err := config.ReadEnv(envFilePath, &c)
	if err != nil {
		return nil, fmt.Errorf("reading env: %w", err)
	}

	return &c, nil
}

func New(ctx context.Context, log logrus.Entry) (App, error) {
	// Build postgres pool
	log.Infoln("Creating pgxpool")

	c, err := ReadEnv(&log)
	if err != nil {
		log.Fatalln(err)
	}

	pgPool, err := postgresql.NewPGXPool(
		ctx,
		&postgresql.ClientConfig{
			Username:     c.PostgresUsername,
			Password:     c.PostgresPassword,
			Host:         c.PostgresHost,
			Port:         c.PostgresPort,
			DatabaseName: c.PostgresDBName,
		},
		log.WithField("loc", "pgxpool"),
	)
	if err != nil {
		log.Errorf("error creating pgxpool: %s\n", err)
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	})

	clickhouseConn, err := pkgClickHouse.Connect(
		pkgClickHouse.Cfg{
			Addr:     c.ClickhouseAddr,
			DB:       c.ClickhouseDBName,
			Username: c.ClickhouseUsername,
			Password: c.ClickhousePassword,
		},
	)
	if err != nil {
		log.Errorf("error creating clickhouse conn: %s\n", err)
		return nil, err
	}

	storages := product.NewStorages(
		pgPool,
		redisClient,
		c.RedisExpirationSeconds,
		clickhouseConn,
	)

	nc, err := nats.Connect(c.NatsURL)
	if err != nil {
		log.Errorf("error creating nats conn: %s\n", err)
		return nil, err
	}

	messageBrokerLogger := log.WithField("loc", "messageBrokerErrTube")
	topicsErrTube := func(err error) {
		messageBrokerLogger.Errorf("topicsErrTube - caught: %s", err)
	}

	goodsChangeQueue := natsPkg.NewQueue(
		*natsPkg.NewPubQueue[event.GoodsChange](nc, c.GoodsChangeKey),
		*natsPkg.NewSubQueue[event.GoodsChange](nc, topicsErrTube, c.GoodsChangeKey, c.NatsFetchTime),
	)

	services, err := product.NewServices(
		*log.WithField("layer", "service"),
		storages,
		*goodsChangeQueue,
	)

	if err != nil {
		log.Errorf("error creating services: %s\n", err)
		return nil, err
	}

	controllers := product.NewControllers(
		*log.WithField("layer", "controller"),
		*services,
		c.APIHost,
		c.APIPort,
		*goodsChangeQueue,
	)

	return app{
		logEntry:    log,
		controllers: controllers,
	}, nil
}

func (a app) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	a.logEntry.Infoln("Go controller api")
	grp.Go(func() error {
		return a.controllers.API.Run(ctx)
	})

	a.logEntry.Infoln("Go monitor")
	grp.Go(func() error {
		return a.controllers.Monitor.Run(ctx)
	})

	return grp.Wait()
}
