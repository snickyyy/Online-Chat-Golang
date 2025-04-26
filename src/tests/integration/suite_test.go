package integration

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"libs/src/settings"
	"libs/src/settings/server"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
)

type AppTestSuite struct {
	suite.Suite
	Ctx    context.Context
	client *http.Client
}

func (suite *AppTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	settings.AppVar = &settings.App{
		Ctx:    ctx,
		Cancel: cancel,
	}
	baseCfg := GetTestConfig()

	db, err := settings.GetDb(baseCfg)
	if err != nil {
		suite.FailNow("Failed to initialize db", err)
	}

	logger, err := settings.GetLogger(baseCfg)
	if err != nil {
		suite.FailNow("Failed to initialize logger", err)
	}

	mongo, err := settings.GetMongoClient(baseCfg, "Testing")
	if err != nil {
		suite.FailNow("Failed to initialize MongoDB client", err)
	}

	redisDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", baseCfg.RedisConfig.Host, baseCfg.RedisConfig.Port),
		Password: baseCfg.RedisConfig.Password,
		DB:       1,
		Protocol: 2,
	})

	mail := gomail.Dialer{}

	app := settings.NewApp(
		db,
		logger,
		baseCfg,
		mongo,
		redisDB,
		&mail,
	)
	settings.AppVar = app
	settings.MakeMigrations(settings.AppVar)

	go func() {
		server.RunServer()
	}()

	time.Sleep(500 * time.Millisecond)

	suite.client = &http.Client{}
	suite.Ctx = ctx
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
