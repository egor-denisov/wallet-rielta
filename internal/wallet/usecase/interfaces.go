package usecase

import (
	"context"

	"github.com/egor-denisov/wallet-rielta/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type (
	Wallet interface {
		CreateNewWalletWithDefaultBalance(ctx context.Context) (*entity.Wallet, error)
		SendFunds(ctx context.Context, from string, to string, amount uint) error
		GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error)
		GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error)
	}

	WalletGateway interface {
		CreateNewWalletWithBalance(ctx context.Context, balance uint) (*entity.Wallet, error)
		SendFunds(ctx context.Context, from string, to string, amount uint) error
		GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error)
		GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error)
	}
)
