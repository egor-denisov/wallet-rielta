package app

import (
	"log/slog"

	"github.com/egor-denisov/wallet-rielta/config"
	v1 "github.com/egor-denisov/wallet-rielta/internal/wallet/controller/http/v1"
	gateway "github.com/egor-denisov/wallet-rielta/internal/wallet/gateway/rabbitmq"
	walletUC "github.com/egor-denisov/wallet-rielta/internal/wallet/usecase"
	amqprpc "github.com/egor-denisov/wallet-rielta/internal/walletWorker/controller/amqp_rpc"
	repo "github.com/egor-denisov/wallet-rielta/internal/walletWorker/repository/postgres"
	workerUC "github.com/egor-denisov/wallet-rielta/internal/walletWorker/usecase"
	"github.com/egor-denisov/wallet-rielta/pkg/httpserver"
	"github.com/egor-denisov/wallet-rielta/pkg/postgres"
	rmqclient "github.com/egor-denisov/wallet-rielta/pkg/rabbitmq/rmq_rpc/client"
	rmqserver "github.com/egor-denisov/wallet-rielta/pkg/rabbitmq/rmq_rpc/server"
	"github.com/gin-gonic/gin"
)

type App struct {
	HTTPServer *httpserver.Server
	RMQServer  *rmqserver.Server
	DB         *postgres.Postgres
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Connect postgres db
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic("app - Run - postgres.New: " + err.Error())
	}
	// Migrate database schema
	err = pg.Migrate("./migrations/20240418133357_init.up.sql")
	if err != nil {
		panic("app - Run - pg.Migrate: " + err.Error())
	}
	// Connect to rabbitmq
	rmqClient, err := rmqclient.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange)
	if err != nil {
		panic("app - Run - rmqServer - server.New" + err.Error())
	}

	// Use cases
	walletUseCase := walletUC.NewWallet(
		gateway.New(rmqClient),
		walletUC.Timeout(cfg.App.Timeout),
		walletUC.DefaultBalance(cfg.App.DefaultBalance),
	)

	workerUseCase := workerUC.NewWalletWorker(
		repo.New(pg),
	)
	// Init http server
	handler := gin.New()
	v1.NewRouter(handler, log, walletUseCase)
	httpServer := httpserver.New(log, handler, httpserver.Port(cfg.HTTP.Port), httpserver.WriteTimeout(cfg.HTTP.Timeout))

	// Init rabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(workerUseCase)

	rmqServer, err := rmqserver.New(
		cfg.RMQ.URL,
		cfg.RMQ.ServerExchange,
		rmqRouter,
		log,
		rmqserver.DefaultGoroutinesCount(cfg.App.CountWorkers),
	)
	if err != nil {
		panic("app - Run - rmqServer - server.New" + err.Error())
	}

	return &App{
		HTTPServer: httpServer,
		RMQServer:  rmqServer,
		DB:         pg,
	}
}
