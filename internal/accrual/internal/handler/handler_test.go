package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mock_storage "github.com/kripsy/gophermart/internal/accrual/internal/mocks"
	"github.com/kripsy/gophermart/internal/accrual/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
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
			req := httptest.NewRequest(tt.method, tt.url, nil)

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
		})

	}
}

func TestHandler_ReadOrdersHandler(t *testing.T) {
	h := &Handler{
		ctx: context.Background(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_storage.NewMockStore(mockCtrl)

	create_time := pgtype.Timestamptz{
		Time:             time.Now(),
		InfinityModifier: 1,
	}
	mockStore.EXPECT().PutOrder(h.ctx, 105626471848586).Return(models.Order{ID: 1, Number: 105626471848586, Status: "PROCESSED", Accrual: 586, UploadedAt: create_time, ProcessedAt: create_time}, nil).AnyTimes()

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
			url:                 "http://127.0.0.1:8080//api/orders/1234567890", //Не могу в хендлере получить номер заказа
			body:                nil,
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

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
		})
	}
}

func TestHandlerCreateOrderHandler(t *testing.T) {
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
	mockStore.EXPECT().PutOrder(gomock.Any(), int64(12351)).Return(models.Order{ID: 1, Number: 12351, Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: nowTime}, nil).AnyTimes()
	mockStore.EXPECT().PutOrder(gomock.Any(), int64(12369)).Return(models.Order{ID: 1, Number: 12369, Status: "PROCESSED", Accrual: 586, UploadedAt: nowTime, ProcessedAt: zeroTime}, nil).AnyTimes()

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
			name:                "CreateOrderHandler 400 no body",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080/api/orders",
			body:                nil,
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 400 wrong number",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080//api/orders",
			body:                []byte(`{"order": "a12344", "goods": [{"description": "Стиральная машинка LG", "price": 47399}, {"description": "Телевизор DLsnXYotclMT", "price": 14599}]}`),
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 422 Luhn is invalid",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080//api/orders",
			body:                []byte(`{"order": "12345", "goods": [{"description": "Стиральная машинка LG", "price": 47399}, {"description": "Телевизор DLsnXYotclMT", "price": 14599}]}`),
			expectedCode:        http.StatusUnprocessableEntity,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 409",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080//api/orders",
			body:                []byte(`{"order": "12351", "goods": [{"description": "Стиральная машинка LG", "price": 47399}, {"description": "Телевизор DLsnXYotclMT", "price": 14599}]}`),
			expectedCode:        http.StatusConflict,
			expectedBody:        "",
			expectedContentType: "",
		},
		{
			name:                "CreateOrderHandler 202",
			method:              http.MethodPost,
			url:                 "http://127.0.0.1:8080//api/orders",
			body:                []byte(`{"order": "12369", "goods": [{"description": "Стиральная машинка LG", "price": 47399}, {"description": "Телевизор DLsnXYotclMT", "price": 14599}]}`),
			expectedCode:        http.StatusAccepted,
			expectedBody:        "",
			expectedContentType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))

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
		})
	}
}
