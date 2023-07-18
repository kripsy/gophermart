package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type RequestOrder struct {
	Number string         `json:"order"`
	Goods  []RequestGoods `json:"goods"`
	//Goods int `json:"goods"`
}

type RequestGoods struct {
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}

type Order struct {
	ID          int64
	Number      int64
	Status      string
	Accrual     int
	UploadedAt  pgtype.Timestamptz
	ProcessedAt pgtype.Timestamptz
}

type ResponseOrder struct {
	Number  string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

const (
	StatusRegistered = "REGISTERED" // Заказ загружен в систему, но не попал в обработку;
	StatusProcessing = "PROCESSING" // Вознаграждение за заказ рассчитывается;
	StatusInvalid    = "INVALID"    // Система расчёта вознаграждений отказала в расчёте;
	StatusProcessed  = "PROCESSED"  // Данные по заказу проверены и информация о расчёте успешно получена.
)
