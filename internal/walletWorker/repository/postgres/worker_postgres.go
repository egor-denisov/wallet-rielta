package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/egor-denisov/wallet-rielta/internal/entity"
	"github.com/egor-denisov/wallet-rielta/pkg/postgres"
)

// WalletRepo -.
type WalletRepo struct {
	*postgres.Postgres
}

// NewWalletRepo -.
func New(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{pg}
}

// CreateNewWallet - creating new wallet entry  in the db.
func (r *WalletRepo) CreateNewWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	_, err := r.DB.ModelContext(ctx, wallet).
		Insert()

	if err != nil {
		return nil, fmt.Errorf("WalletRepo - CreateNewWallet - r.DB: %w", err)
	}

	return wallet, nil
}

// SendFunds - decreasing the balance of the sender and an increasing the receiver.
// Adding an entry to a transaction table.
func (r *WalletRepo) SendFunds(ctx context.Context, transaction *entity.Transaction) error {
	// Using the db transaction
	tx, err := r.DB.BeginContext(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}

	defer tx.Close()
	// Decreasing the balance of the sender
	res, err := r.DB.ModelContext(ctx, new(entity.Wallet)).
		Set("balance = balance - ?", transaction.Amount).
		Where("id = ?", transaction.From).
		Update()
	// If error then rollback the transaction
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	// If walletId is not found then return 404
	if res.RowsAffected() == 0 {
		_ = tx.Rollback()
		return entity.ErrWalletNotFound
	}
	// Increasing the balance of the receiver
	res, err = r.DB.ModelContext(ctx, new(entity.Wallet)).
		Set("balance = balance + ?", transaction.Amount).
		Where("id = ?", transaction.To).
		Update()
	// If error or walletId is not found then rollback the transaction
	if err != nil || res.RowsAffected() == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	// Adding an entry to a transaction table
	_, err = r.DB.ModelContext(ctx, transaction).
		Insert()

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	// Make commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}

	return nil
}

// GetWalletHistoryByID - getting all transaction records from the user with the walletID.
func (r *WalletRepo) GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	err := r.DB.ModelContext(ctx, new(entity.Wallet)).
		Where("id = ?", walletID).
		Select()

	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, entity.ErrWalletNotFound
		}

		return nil, fmt.Errorf("WalletRepo - GetWalletHistoryByID - r.DB: %w", err)
	}

	err = r.DB.ModelContext(ctx, &transactions).
		Where("from_wallet_id = ?", walletID).
		WhereOr("to_wallet_id = ?", walletID).
		Select()

	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletHistoryByID - r.DB: %w", err)
	}

	return transactions, nil
}

// GetWalletByID - getting wallet info by walletID.
func (r *WalletRepo) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	wallet := new(entity.Wallet)
	if walletID == "test" {
		time.Sleep(4 * time.Second)
	}

	err := r.DB.ModelContext(ctx, wallet).
		Where("id = ?", walletID).
		Select()

	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, entity.ErrWalletNotFound
		}

		return nil, fmt.Errorf("WalletRepo - GetWalletByID - r.DB: %w", err)
	}

	return wallet, nil
}
