package entity

import (
	"context"
	"errors"

	rmqrpc "github.com/egor-denisov/wallet-rielta/pkg/rabbitmq/rmq_rpc"
)

var (
	// Wallet errors.
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrWrongAmount      = errors.New("wrong amount")
	ErrSenderIsReceiver = errors.New("sender is receiver")
	ErrEmptyWallet      = errors.New("wallet address is empty")

	// Requset errors.
	ErrTimeout  = context.DeadlineExceeded
	ErrNotFound = rmqrpc.ErrNotFound
)
