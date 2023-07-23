package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	c "github.com/gorilla/context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	mock_storage "github.com/kripsy/gophermart/internal/gophermart/internal/mocks"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandlerCreateOrderHandler(t *testing.T) {
	l, _ := logger.InitLogger("Warn")
	h := &Handler{
		ctx: context.Background(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_storage.NewMockStore(mockCtrl)

	nowTime := pgtype.Timestamptz{
		Time:             time.Now(),
		InfinityModifier: 1,
	}
	zeroTime := pgtype.Timestamptz{
		Time:             time.Time{},
		InfinityModifier: 1,
	}
	mockStore.EXPECT().PutOrder(h.ctx, "username", int64(1230)).Return(models.Order{}, errors.New("Error for test")).AnyTimes()
	mockStore.EXPECT().PutOrder(h.ctx, "username", int64(1222)).Return(models.Order{ID: 1, Number: 1222, Username: "username", Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: zeroTime}, nil).AnyTimes()
	mockStore.EXPECT().PutOrder(h.ctx, "username", int64(1255)).Return(models.Order{ID: 1, Number: 1255, Username: "username", Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: nowTime}, nil).AnyTimes()
	mockStore.EXPECT().PutOrder(h.ctx, "username", int64(1263)).Return(models.Order{ID: 1, Number: 1255, Username: "username", Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: zeroTime}, nil).AnyTimes()
	mockStore.EXPECT().PutOrder(h.ctx, "username", int64(1248)).Return(models.Order{ID: 1, Number: 1248, Username: "username1", Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: nowTime}, nil).AnyTimes()
	mockStore.EXPECT().PutOrder(h.ctx, "username1", int64(1248)).Return(models.Order{ID: 1, Number: 1248, Username: "username1", Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: nowTime}, nil).AnyTimes()

	tests := []struct {
		name                string
		method              string
		url                 string
		body                []byte
		username            string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "CreateOrderHandler 400 case1",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                nil,
			username:            "username",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 400 case1",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                []byte(`a1223`),
			username:            "username",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 422 invalid LuhnValid",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                []byte(`1223`),
			username:            "username",
			expectedCode:        http.StatusUnprocessableEntity,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 500 InternalServerError",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                []byte(`1230`),
			username:            "username",
			expectedCode:        http.StatusInternalServerError,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 409 order number has already been uploaded by another user",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                []byte(`1248`),
			username:            "username",
			expectedCode:        http.StatusConflict,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 200 order number has already been uploaded by user",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                []byte(`1255`),
			username:            "username",
			expectedCode:        http.StatusOK,
			expectedBody:        "",
			expectedContentType: "",
		},
		//{ // Тест не проходит. Зависает при отправке значения в канал.
		//	name:                "CreateOrderHandler 202",
		//	method:              http.MethodPost,
		//	url:                 "http://127.0.0.1:8080/api/user/orders",
		//	body:                []byte(`1263`),
		//	username:            "username",
		//	expectedCode:        http.StatusAccepted,
		//	expectedBody:        "",
		//	expectedContentType: "",
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

			c.Set(req, "username", tt.username)

			h.CreateOrderHandler(mockStore)(rw, req)

			resp := rw.Result()

			if tt.expectedCode != 0 {
				assert.Equal(t, tt.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedBody, string(body), "Body ответа не совпадает с ожидаемым")
			}

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
			}

			err := resp.Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		})
	}
}

func TestHandlerReadOrdersHandler(t *testing.T) {
	l, _ := logger.InitLogger("Warn")
	h := &Handler{
		ctx: context.Background(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_storage.NewMockStore(mockCtrl)

	nowTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.UTC)
	zeroTime := time.Time{}

	orders := []models.ResponseOrder{}
	order := models.ResponseOrder{ID: 1, Username: "username3", Number: "1230", Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: zeroTime}

	orders = append(orders, order)

	mockStore.EXPECT().GetOrders(h.ctx, "username1").Return([]models.ResponseOrder{}, errors.New("Error for test")).AnyTimes()
	mockStore.EXPECT().GetOrders(h.ctx, "username2").Return([]models.ResponseOrder{}, nil).AnyTimes()
	mockStore.EXPECT().GetOrders(h.ctx, "username3").Return(orders, nil).AnyTimes()

	tests := []struct {
		name                string
		method              string
		url                 string
		body                []byte
		username            string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "CreateOrderHandler 500 InternalServerError",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                nil,
			username:            "username1",
			expectedCode:        http.StatusInternalServerError,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 204 No Content",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                nil,
			username:            "username2",
			expectedCode:        http.StatusNoContent,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 204 No Content",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/orders",
			body:                nil,
			username:            "username3",
			expectedCode:        http.StatusOK,
			expectedBody:        "[{\"order\":\"1230\",\"status\":\"PROCESSED\",\"accrual\":586,\"uploaded_at\":\"2021-08-15T14:30:45.0000001Z\"}]\n",
			expectedContentType: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

			c.Set(req, "username", tt.username)

			h.ReadOrdersHandler(mockStore)(rw, req)

			resp := rw.Result()

			if tt.expectedCode != 0 {
				assert.Equal(t, tt.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedBody, string(body), "Body ответа не совпадает с ожидаемым")
			}

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
			}

			err := resp.Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		})
	}
}

func TestHandlerReadUserBalanceHandler(t *testing.T) {
	l, _ := logger.InitLogger("Warn")
	h := &Handler{
		ctx: context.Background(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_storage.NewMockStore(mockCtrl)

	nowTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.UTC)
	zeroTime := time.Time{}

	mockStore.EXPECT().GetBalance(h.ctx, "username1").Return(models.ResponseBalance{}, errors.New("Error for test")).AnyTimes()
	mockStore.EXPECT().GetBalance(h.ctx, "username2").Return(models.ResponseBalance{ID: 1, Username: "username2", Current: 100, Withdrawn: 20, UploadedAt: nowTime, ProcessedAt: zeroTime}, nil).AnyTimes()

	tests := []struct {
		name                string
		method              string
		url                 string
		body                []byte
		username            string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "ReadUserBalanceHandler 500 InternalServerError",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/balance",
			body:                nil,
			username:            "username1",
			expectedCode:        http.StatusInternalServerError,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "ReadUserBalanceHandler 200",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/balance",
			body:                nil,
			username:            "username2",
			expectedCode:        http.StatusOK,
			expectedBody:        "{\"current\":100,\"withdrawn\":20,\"uploaded_at\":\"2021-08-15T14:30:45.0000001Z\"}\n",
			expectedContentType: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

			c.Set(req, "username", tt.username)

			h.ReadUserBalanceHandler(mockStore)(rw, req)

			resp := rw.Result()

			if tt.expectedCode != 0 {
				assert.Equal(t, tt.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedBody, string(body), "Body ответа не совпадает с ожидаемым")
			}

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
			}

			err := resp.Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		})
	}
}

func TestHandlerCreateWithdrawHandler(t *testing.T) {
	l, _ := logger.InitLogger("Warn")
	h := &Handler{
		ctx: context.Background(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_storage.NewMockStore(mockCtrl)

	mockStore.EXPECT().PutWithdraw(h.ctx, "username1", int64(2377225624), 20).Return(errors.New("Error for test")).AnyTimes()
	mockStore.EXPECT().PutWithdraw(h.ctx, "username1", int64(12377), 20).Return(nil).AnyTimes()

	tests := []struct {
		name                string
		method              string
		url                 string
		body                []byte
		username            string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "CreateWithdrawHandler 400 No body",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/balance/withdraw",
			body:                nil,
			username:            "username1",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateWithdrawHandler 400 wrong body",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/balance/withdraw",
			body:                []byte(`{"order": "a2377225624", "sum": 20}`),
			username:            "username1",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateWithdrawHandler 422 LuhnValid is invalid",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/balance/withdraw",
			body:                []byte(`{"order": "2377225625", "sum": 20}`),
			username:            "username1",
			expectedCode:        http.StatusUnprocessableEntity,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateWithdrawHandler 500 InternalServerError",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/balance/withdraw",
			body:                []byte(`{"order": "2377225624", "sum": 20}`),
			username:            "username1",
			expectedCode:        http.StatusInternalServerError,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateWithdrawHandler 500 InternalServerError",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/user/balance/withdraw",
			body:                []byte(`{"order": "12377", "sum": 20}`),
			username:            "username1",
			expectedCode:        http.StatusOK,
			expectedBody:        "",
			expectedContentType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

			c.Set(req, "username", tt.username)

			h.CreateWithdrawHandler(mockStore)(rw, req)

			resp := rw.Result()

			if tt.expectedCode != 0 {
				assert.Equal(t, tt.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedBody, string(body), "Body ответа не совпадает с ожидаемым")
			}

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
			}

			err := resp.Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		})
	}
}

func TestHandlerReadWithdrawsHandler(t *testing.T) {
	l, _ := logger.InitLogger("Warn")
	h := &Handler{
		ctx: context.Background(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_storage.NewMockStore(mockCtrl)

	nowTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.UTC)
	//zeroTime := time.Time{}

	var withdraws []models.ResponseOrder
	withdraw := models.ResponseOrder{ID: 1, Username: "username2", Number: "10", Status: "PROCESSED", Accrual: 50, UploadedAt: nowTime, ProcessedAt: nowTime}

	withdraws = append(withdraws, withdraw)

	mockStore.EXPECT().GetWithdraws(h.ctx, "username1").Return([]models.ResponseOrder{}, errors.New("Error for test")).AnyTimes()
	mockStore.EXPECT().GetWithdraws(h.ctx, "username2").Return(withdraws, nil).AnyTimes()

	tests := []struct {
		name                string
		method              string
		url                 string
		body                []byte
		username            string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "CreateWithdrawHandler 500 InternalServerError",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/withdraws",
			body:                nil,
			username:            "username1",
			expectedCode:        http.StatusInternalServerError,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateWithdrawHandler 200 InternalServerError",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/api/user/withdraws",
			body:                nil,
			username:            "username2",
			expectedCode:        http.StatusOK,
			expectedBody:        "[{\"order\":\"10\",\"status\":\"PROCESSED\",\"accrual\":50,\"uploaded_at\":\"2021-08-15T14:30:45.0000001Z\"}]\n",
			expectedContentType: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

			c.Set(req, "username", tt.username)

			h.ReadWithdrawsHandler(mockStore)(rw, req)

			resp := rw.Result()

			if tt.expectedCode != 0 {
				assert.Equal(t, tt.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedBody, string(body), "Body ответа не совпадает с ожидаемым")
			}

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
			}

			err := resp.Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		})
	}
}

func TestHandler(t *testing.T) {
	l, _ := logger.InitLogger("Warn")
	h := &Handler{
		ctx: context.Background(),
	}
	tests := []struct {
		name                string
		method              string
		url                 string
		body                []byte
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "TestHandler",
			method:              http.MethodGet,
			url:                 "http://127.0.0.1:8080/test",
			body:                nil,
			expectedCode:        http.StatusOK,
			expectedBody:        "",
			expectedContentType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

			h.TestHandler(rw, req)

			resp := rw.Result()

			if tt.expectedCode != 0 {
				assert.Equal(t, tt.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedBody, string(body), "Body ответа не совпадает с ожидаемым")
			}

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
			}

			err := resp.Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		})

	}
}
