package amqprpc

import (
	"github.com/egor-denisov/wallet-rielta/internal/walletWorker/usecase"
	"github.com/egor-denisov/wallet-rielta/pkg/rabbitmq/rmq_rpc/server"
)

func NewRouter(r usecase.WalletWorker) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newWalletWorkerRoutes(routes, r)
	}

	return routes
}
