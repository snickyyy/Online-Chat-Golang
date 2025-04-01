package integration

import (
	"libs/src/settings"
	"libs/src/settings/server"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
)

type AppTestSuite struct {
	suite.Suite
	client *http.Client
}

func (suite *AppTestSuite) SetupSuite() {
	settings.InitContext()
	baseCfg := GetTestConfig()

	pgPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	postgresBaseDb := settings.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Port:     pgPort,
		Sslmode:  os.Getenv("DB_SSL_MODE"),
	}
	db := SetupTestDatabase(postgresBaseDb, baseCfg.PostgresConfig)

	logger, err := settings.GetLogger(baseCfg)
	if err != nil {
		suite.FailNow("Failed to initialize logger", err)
	}

	mongo, err := settings.GetMongoClient(baseCfg, "Testing")
	if err != nil {
		suite.FailNow("Failed to initialize MongoDB client", err)
	}

	redis := settings.NewRedisClient(baseCfg)
	mail := gomail.Dialer{}

	app := settings.NewApp(
		db,
		logger,
		baseCfg,
		mongo,
		redis,
		&mail,
	)
	settings.AppVar = app
	settings.MakeMigrations(settings.AppVar)

	go func() {
		server.RunServer()
	}()

	time.Sleep(500 * time.Millisecond)

	suite.client = &http.Client{}
}

func (suite *AppTestSuite) TearDownSuite() {
	pgPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	postgresBaseDb := settings.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Port:     pgPort,
		Sslmode:  os.Getenv("DB_SSL_MODE"),
	}
	DropTestDatabase(postgresBaseDb)

	settings.AppVar.MongoDB.Drop(settings.Context.Ctx)

	settings.AppVar.RedisClient.FlushAll(settings.Context.Ctx)
	settings.AppVar.RedisClient.Close()
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
