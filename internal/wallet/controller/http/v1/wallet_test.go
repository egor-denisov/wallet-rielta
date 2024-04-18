package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"

	"github.com/egor-denisov/wallet-rielta/internal/entity"
	mock_usecase "github.com/egor-denisov/wallet-rielta/internal/wallet/usecase/mocks"
	"github.com/egor-denisov/wallet-rielta/pkg/logger"
)

var errSomethingWrong = errors.New("something went wrong")

func Test_createNewWallet(t *testing.T) {
	for _, test := range testsCreateNewWallet {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id)
			handler := walletRoutes{
				w: repo,
				l: logger.SetupLogger("debug"),
			}
			// Init Endpoint
			r := gin.New()
			r.POST("/", handler.createNewWallet)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

var testsCreateNewWallet = []struct {
	name                 string
	mockBehavior         func(r *mock_usecase.MockWallet, id string)
	id                   string
	expectedStatusCode   int
	expectedResponseBody string
}{
	{
		name: "Ok",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().CreateNewWalletWithDefaultBalance(context.Background()).Return(&entity.Wallet{
				ID:      id,
				Balance: 100,
			}, nil)
		},
		id:                   "5b53700ed469fa6a09ea72bb78f36fd9",
		expectedStatusCode:   200,
		expectedResponseBody: `{"id":"5b53700ed469fa6a09ea72bb78f36fd9","balance":100}`,
	},
	{
		name: "Something went wrong",
		mockBehavior: func(r *mock_usecase.MockWallet, _ string) {
			r.EXPECT().CreateNewWalletWithDefaultBalance(context.Background()).Return(nil, errSomethingWrong)
		},
		expectedStatusCode:   500,
		expectedResponseBody: "",
	},
	{
		name: "Timeout",
		mockBehavior: func(r *mock_usecase.MockWallet, _ string) {
			r.EXPECT().CreateNewWalletWithDefaultBalance(context.Background()).Return(nil, entity.ErrTimeout)
		},
		expectedStatusCode:   504,
		expectedResponseBody: "",
	},
}

func Test_sendFunds(t *testing.T) {
	for _, test := range testSendFunds {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id, test.req)
			handler := walletRoutes{
				w: repo,
				l: logger.SetupLogger("debug"),
			}
			// Init Endpoint
			r := gin.New()
			r.POST("/:walletId/send", handler.sendFunds)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/send", test.id), bytes.NewBufferString(test.reqBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

var testSendFunds = []struct {
	name                 string
	id                   string
	reqBody              string
	req                  transactionRequest
	mockBehavior         func(r *mock_usecase.MockWallet, id string, req transactionRequest)
	expectedStatusCode   int
	expectedResponseBody string
}{
	{
		name:    "Ok",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"to":"eb376add88bf8e70f80787266a0801d5","amount":100}`,
		req: transactionRequest{
			To:     "eb376add88bf8e70f80787266a0801d5",
			Amount: 100,
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(nil)
		},
		expectedStatusCode:   200,
		expectedResponseBody: "",
	},
	{
		name:    "Not found",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"to":"eb376add88bf8e70f80787266a0801d5","amount":100}`,
		req: transactionRequest{
			To:     "eb376add88bf8e70f80787266a0801d5",
			Amount: 100,
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(entity.ErrWalletNotFound)
		},
		expectedStatusCode:   404,
		expectedResponseBody: "",
	},
	{
		name:    "Wrong input - without receiver id",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"amount":100}`,
		req: transactionRequest{
			Amount: 100,
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(entity.ErrEmptyWallet)
		},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:    "Wrong input - without amount",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"to":"eb376add88bf8e70f80787266a0801d5"}`,
		req: transactionRequest{
			To: "eb376add88bf8e70f80787266a0801d5",
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(entity.ErrWrongAmount)
		},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:                 "Wrong input - amount less 0",
		id:                   "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody:              `{"to":"eb376add88bf8e70f80787266a0801d5","amount":-100}`,
		mockBehavior:         func(_ *mock_usecase.MockWallet, _ string, _ transactionRequest) {},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:                 "Wrong input - amount not number",
		id:                   "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody:              `{"to":"eb376add88bf8e70f80787266a0801d5","amount":"abc"}`,
		mockBehavior:         func(_ *mock_usecase.MockWallet, _ string, _ transactionRequest) {},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:                 "Wrong input - not json",
		id:                   "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody:              `helloworld`,
		mockBehavior:         func(_ *mock_usecase.MockWallet, _ string, _ transactionRequest) {},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:    "Wrong input - empty request body",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{}`,
		req:     transactionRequest{},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(entity.ErrWrongAmount)
		},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:    "Sender is receiver",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"to":"5b53700ed469fa6a09ea72bb78f36fd9","amount":100}`,
		req: transactionRequest{
			To:     "5b53700ed469fa6a09ea72bb78f36fd9",
			Amount: 100,
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(entity.ErrSenderIsReceiver)
		},
		expectedStatusCode:   400,
		expectedResponseBody: "",
	},
	{
		name:    "Timeout",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"to":"eb376add88bf8e70f80787266a0801d5","amount":100}`,
		req: transactionRequest{
			To:     "eb376add88bf8e70f80787266a0801d5",
			Amount: 100,
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(entity.ErrTimeout)
		},
		expectedStatusCode:   504,
		expectedResponseBody: "",
	},
	{
		name:    "Something went wrong",
		id:      "5b53700ed469fa6a09ea72bb78f36fd9",
		reqBody: `{"to":"eb376add88bf8e70f80787266a0801d5","amount":100}`,
		req: transactionRequest{
			To:     "eb376add88bf8e70f80787266a0801d5",
			Amount: 100,
		},
		mockBehavior: func(r *mock_usecase.MockWallet, id string, req transactionRequest) {
			r.EXPECT().SendFunds(context.Background(), id, req.To, req.Amount).Return(errSomethingWrong)
		},
		expectedStatusCode:   500,
		expectedResponseBody: "",
	},
}

func Test_getWalletHistoryByID(t *testing.T) {
	for _, test := range testGetWalletHistoryByID {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id)
			handler := walletRoutes{
				w: repo,
				l: logger.SetupLogger("debug"),
			}
			// Init Endpoint
			r := gin.New()
			r.GET("/:walletId/history", handler.GetWalletHistoryByID)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/history", test.id), nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

var testGetWalletHistoryByID = []struct {
	name                 string
	id                   string
	mockBehavior         func(r *mock_usecase.MockWallet, id string)
	expectedStatusCode   int
	expectedResponseBody string
}{
	{
		name: "Ok - history exists (sending and receiving)",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			t, _ := time.Parse(time.RFC3339, "2024-02-04T17:25:35.448Z")

			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return([]entity.Transaction{
				{
					Time:   t,
					From:   "5b53700ed469fa6a09ea72bb78f36fd9",
					To:     "eb376add88bf8e70f80787266a0801d5",
					Amount: 30,
				},
				{
					Time:   t,
					From:   "eb376add88bf8e70f80787266a0801d5",
					To:     "5b53700ed469fa6a09ea72bb78f36fd9",
					Amount: 30,
				},
			}, nil)
		},
		expectedStatusCode: 200,
		expectedResponseBody: `[{"time":"2024-02-04T17:25:35.448Z","from":"5b53700ed469fa6a09ea72bb78f36fd9",` +
			`"to":"eb376add88bf8e70f80787266a0801d5","amount":30},` +
			`{"time":"2024-02-04T17:25:35.448Z","from":"eb376add88bf8e70f80787266a0801d5",` +
			`"to":"5b53700ed469fa6a09ea72bb78f36fd9","amount":30}]`,
	},
	{
		name: "Ok - history exists (only sending)",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			t, _ := time.Parse(time.RFC3339, "2024-02-04T17:25:35.448Z")

			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return([]entity.Transaction{
				{
					Time:   t,
					From:   "5b53700ed469fa6a09ea72bb78f36fd9",
					To:     "eb376add88bf8e70f80787266a0801d5",
					Amount: 30.0,
				},
			}, nil)
		},
		expectedStatusCode: 200,
		expectedResponseBody: `[{"time":"2024-02-04T17:25:35.448Z","from":"5b53700ed469fa6a09ea72bb78f36fd9",` +
			`"to":"eb376add88bf8e70f80787266a0801d5","amount":30}]`,
	},
	{
		name: "Ok - history exists (only receiving)",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			t, _ := time.Parse(time.RFC3339, "2024-02-04T17:25:35.448Z")

			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return([]entity.Transaction{
				{
					Time:   t,
					From:   "eb376add88bf8e70f80787266a0801d5",
					To:     "5b53700ed469fa6a09ea72bb78f36fd9",
					Amount: 30.0,
				},
			}, nil)
		},
		expectedStatusCode: 200,
		expectedResponseBody: `[{"time":"2024-02-04T17:25:35.448Z","from":"eb376add88bf8e70f80787266a0801d5",` +
			`"to":"5b53700ed469fa6a09ea72bb78f36fd9","amount":30}]`,
	},
	{
		name: "Ok - history is empty",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return([]entity.Transaction{}, nil)
		},
		expectedStatusCode:   200,
		expectedResponseBody: `[]`,
	},
	{
		name: "Not Found",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return(nil, entity.ErrWalletNotFound)
		},
		expectedStatusCode:   404,
		expectedResponseBody: "",
	},
	{
		name: "Timeout",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return(nil, entity.ErrTimeout)
		},
		expectedStatusCode:   504,
		expectedResponseBody: "",
	},
	{
		name: "Something went wrong",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletHistoryByID(context.Background(), id).Return(nil, errSomethingWrong)
		},
		expectedStatusCode:   500,
		expectedResponseBody: "",
	},
}

func Test_getWalletByID(t *testing.T) {
	for _, test := range testGetWalletByID {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id)
			handler := walletRoutes{
				w: repo,
				l: logger.SetupLogger("debug"),
			}
			// Init Endpoint
			r := gin.New()
			r.GET("/:walletId", handler.GetWalletByID)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/"+test.id, nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

var testGetWalletByID = []struct {
	name                 string
	id                   string
	mockBehavior         func(r *mock_usecase.MockWallet, id string)
	expectedStatusCode   int
	expectedResponseBody string
}{
	{
		name: "Ok",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletByID(context.Background(), id).Return(&entity.Wallet{
				ID:      id,
				Balance: 100,
			}, nil)
		},
		expectedStatusCode:   200,
		expectedResponseBody: `{"id":"5b53700ed469fa6a09ea72bb78f36fd9","balance":100}`,
	},
	{
		name: "Not Found",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletByID(context.Background(), id).Return(nil, entity.ErrWalletNotFound)
		},
		expectedStatusCode:   404,
		expectedResponseBody: "",
	},
	{
		name: "Timeout",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletByID(context.Background(), id).Return(nil, entity.ErrTimeout)
		},
		expectedStatusCode:   504,
		expectedResponseBody: "",
	},
	{
		name: "Something went wrong",
		id:   "5b53700ed469fa6a09ea72bb78f36fd9",
		mockBehavior: func(r *mock_usecase.MockWallet, id string) {
			r.EXPECT().GetWalletByID(context.Background(), id).Return(nil, errSomethingWrong)
		},
		expectedStatusCode:   500,
		expectedResponseBody: "",
	},
}
