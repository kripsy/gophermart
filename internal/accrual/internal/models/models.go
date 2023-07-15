package models

import "github.com/jackc/pgx/v5/pgtype"

type RequestOrder struct {
	Number string         `json:"order"`
	Goods  []RequestGoods `json:"goods"`
}

type RequestGoods struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type Order struct {
	ID          int64
	UserName    string
	Number      int64
	Status      string
	Accrual     int
	UploadedAt  pgtype.Timestamptz
	ProcessedAt pgtype.Timestamptz
}

const (
	StatusRegistered = "REGISTERED" // Заказ загружен в систему, но не попал в обработку;
	StatusProcessing = "PROCESSING" // Вознаграждение за заказ рассчитывается;
	StatusInvalid    = "INVALID"    // Система расчёта вознаграждений отказала в расчёте;
	StatusProcessed  = "PROCESSED"  // Данные по заказу проверены и информация о расчёте успешно получена.
)
