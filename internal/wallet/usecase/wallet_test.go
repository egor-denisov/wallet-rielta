package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/egor-denisov/wallet-rielta/internal/entity"
	mock_usecase "github.com/egor-denisov/wallet-rielta/internal/wallet/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

var errSomethingWentWrong = errors.New("something went wrong")

func Test_CreateNewWalletWithDefaultBalance(t *testing.T) {
	for _, test := range testsCreateNewWalletWithBalance {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			gateway := mock_usecase.NewMockWalletGateway(c)
			test.mockBehavior(gateway)

			// Call function and check the result
			wallet, err := NewWallet(gateway).CreateNewWalletWithDefaultBalance(context.Background())
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected %v, got %v", test.expectedError, err)
			}

			assert.Equal(t, wallet, test.expectedWallet)
		})
	}
}

var testsCreateNewWalletWithBalance = []struct {
	name           string
	mockBehavior   func(r *mock_usecase.MockWalletGateway)
	expectedError  error
	expectedWallet *entity.Wallet
}{
	{
		name: "Ok",
		mockBehavior: func(r *mock_usecase.MockWalletGateway) {
			r.EXPECT().CreateNewWalletWithBalance(gomock.Any(), _defaultBalance).Return(&entity.Wallet{
				ID:      "5b53700ed469fa6a09ea72bb78f36fd9",
				Balance: 100,
			}, nil)
		},
		expectedError: nil,
		expectedWallet: &entity.Wallet{
			ID:      "5b53700ed469fa6a09ea72bb78f36fd9",
			Balance: 100,
		},
	},
	{
		name: "Something went wrong",
		mockBehavior: func(r *mock_usecase.MockWalletGateway) {
			r.EXPECT().CreateNewWalletWithBalance(gomock.Any(), _defaultBalance).Return(nil, errSomethingWentWrong)
		},
		expectedError:  errSomethingWentWrong,
		expectedWallet: nil,
	},
}

func Test_SendFunds(t *testing.T) {
	for _, test := range testsSendFunds {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			gateway := mock_usecase.NewMockWalletGateway(c)
			test.mockBehavior(gateway, test.from, test.to, test.amount)

			// Call function and check the result
			err := NewWallet(gateway).SendFunds(context.Background(), test.from, test.to, test.amount)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected %v, got %v", test.expectedError, err)
			}
		})
	}
}

var testsSendFunds = []struct {
	name          string
	mockBehavior  func(r *mock_usecase.MockWalletGateway, from, to string, amount uint)
	from          string
	to            string
	amount        uint
	expectedError error
}{
	{
		name: "Ok",
		mockBehavior: func(r *mock_usecase.MockWalletGateway, from, to string, amount uint) {
			r.EXPECT().SendFunds(gomock.Any(), from, to, amount).Return(nil)
		},
		from:          "5b53700ed469fa6a09ea72bb78f36fd9",
		to:            "eb376add88bf8e70f80787266a0801d5",
		amount:        100,
		expectedError: nil,
	},
	{
		name:          "Amount must be greater than 0",
		mockBehavior:  func(_ *mock_usecase.MockWalletGateway, _, _ string, _ uint) {},
		from:          "5b53700ed469fa6a09ea72bb78f36fd9",
		to:            "eb376add88bf8e70f80787266a0801d5",
		amount:        0,
		expectedError: entity.ErrWrongAmount,
	},
	{
		name:          "Wallets ID`s must be non-empty",
		mockBehavior:  func(_ *mock_usecase.MockWalletGateway, _, _ string, _ uint) {},
		from:          "5b53700ed469fa6a09ea72bb78f36fd9",
		to:            "",
		amount:        100,
		expectedError: entity.ErrEmptyWallet,
	},
	{
		name:          "Wallets from and to must be not equal",
		mockBehavior:  func(_ *mock_usecase.MockWalletGateway, _, _ string, _ uint) {},
		from:          "5b53700ed469fa6a09ea72bb78f36fd9",
		to:            "5b53700ed469fa6a09ea72bb78f36fd9",
		amount:        100,
		expectedError: entity.ErrSenderIsReceiver,
	},
	{
		name: "Something went wrong",
		mockBehavior: func(r *mock_usecase.MockWalletGateway, from, to string, amount uint) {
			r.EXPECT().SendFunds(gomock.Any(), from, to, amount).Return(errSomethingWentWrong)
		},
		from:          "5b53700ed469fa6a09ea72bb78f36fd9",
		to:            "eb376add88bf8e70f80787266a0801d5",
		amount:        100,
		expectedError: errSomethingWentWrong,
	},
}

func Test_GetWalletHistoryByID(t *testing.T) {
	for _, test := range testsGetWalletHistoryByID {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			gateway := mock_usecase.NewMockWalletGateway(c)
			test.mockBehavior(gateway, test.walletID)

			// Call function and check the result
			transactions, err := NewWallet(gateway).GetWalletHistoryByID(context.Background(), test.walletID)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected %v, got %v", test.expectedError, err)
			}

			for i, tr := range transactions {
				assert.Equal(t, tr, test.expectedTransactions[i])
			}
		})
	}
}

var testsGetWalletHistoryByID = []struct {
	name                 string
	mockBehavior         func(r *mock_usecase.MockWalletGateway, walletID string)
	walletID             string
	expectedTransactions []entity.Transaction
	expectedError        error
}{
	{
		name: "Ok",
		mockBehavior: func(r *mock_usecase.MockWalletGateway, walletID string) {
			r.EXPECT().GetWalletHistoryByID(gomock.Any(), walletID).Return([]entity.Transaction{
				{
					Time:   time.Date(2024, time.February, 4, 17, 25, 35, 0, time.UTC),
					From:   "5b53700ed469fa6a09ea72bb78f36fd9",
					To:     "eb376add88bf8e70f80787266a0801d5",
					Amount: 30,
				},
				{
					Time:   time.Date(2024, time.February, 4, 17, 25, 35, 0, time.UTC),
					From:   "eb376add88bf8e70f80787266a0801d5",
					To:     "5b53700ed469fa6a09ea72bb78f36fd9",
					Amount: 30,
				},
			}, nil)
		},
		walletID: "5b53700ed469fa6a09ea72bb78f36fd9",
		expectedTransactions: []entity.Transaction{
			{
				Time:   time.Date(2024, time.February, 4, 17, 25, 35, 0, time.UTC),
				From:   "5b53700ed469fa6a09ea72bb78f36fd9",
				To:     "eb376add88bf8e70f80787266a0801d5",
				Amount: 30,
			},
			{
				Time:   time.Date(2024, time.February, 4, 17, 25, 35, 0, time.UTC),
				From:   "eb376add88bf8e70f80787266a0801d5",
				To:     "5b53700ed469fa6a09ea72bb78f36fd9",
				Amount: 30,
			},
		},
		expectedError: nil,
	},
	{
		name: "Something went wrong",
		mockBehavior: func(r *mock_usecase.MockWalletGateway, walletID string) {
			r.EXPECT().GetWalletHistoryByID(gomock.Any(), walletID).Return([]entity.Transaction{}, errSomethingWentWrong)
		},
		walletID:             "5b53700ed469fa6a09ea72bb78f36fd9",
		expectedTransactions: []entity.Transaction{},
		expectedError:        errSomethingWentWrong,
	},
}

func Test_GetWalletByID(t *testing.T) {
	for _, test := range testsGetWalletByID {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			gateway := mock_usecase.NewMockWalletGateway(c)
			test.mockBehavior(gateway, test.walletID)

			// Call function and check the result
			wallet, err := NewWallet(gateway).GetWalletByID(context.Background(), test.walletID)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected %v, got %v", test.expectedError, err)
			}

			assert.Equal(t, wallet, test.expectedWallet)
		})
	}
}

var testsGetWalletByID = []struct {
	name           string
	mockBehavior   func(r *mock_usecase.MockWalletGateway, walletID string)
	walletID       string
	expectedError  error
	expectedWallet *entity.Wallet
}{
	{
		name: "Ok",
		mockBehavior: func(r *mock_usecase.MockWalletGateway, walletID string) {
			r.EXPECT().GetWalletByID(gomock.Any(), walletID).Return(&entity.Wallet{
				ID:      "5b53700ed469fa6a09ea72bb78f36fd9",
				Balance: 100,
			}, nil)
		},
		expectedError: nil,
		expectedWallet: &entity.Wallet{
			ID:      "5b53700ed469fa6a09ea72bb78f36fd9",
			Balance: 100,
		},
	},
	{
		name: "Something went wrong",
		mockBehavior: func(r *mock_usecase.MockWalletGateway, walletID string) {
			r.EXPECT().GetWalletByID(gomock.Any(), walletID).Return(nil, errSomethingWentWrong)
		},
		expectedError:  errSomethingWentWrong,
		expectedWallet: nil,
	},
}
