package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID          int64
	Username    string
	Number      int64
	Status      string
	Accrual     int
	UploadedAt  pgtype.Timestamptz
	ProcessedAt pgtype.Timestamptz
}

type ResponseOrder struct {
	ID          int64     `json:"-"`
	Username    string    `json:"-"`
	Number      string    `json:"order"`
	Status      string    `json:"status"`
	Accrual     int       `json:"accrual,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
	ProcessedAt time.Time `json:"-"`
}

type ResponseWithdrawals struct {
	ID          int64     `json:"-"`
	Username    string    `json:"-"`
	Number      string    `json:"order"`
	Status      string    `json:"-"`
	Accrual     int       `json:"sum,omitempty"`
	UploadedAt  time.Time `json:"-"`
	ProcessedAt time.Time `json:"processed_at"`
}

type ResponseBalance struct {
	ID          int64     `json:"-"`
	Username    string    `json:"-"`
	Current     int       `json:"current"`
	Withdrawn   int       `json:"withdrawn"`
	UploadedAt  time.Time `json:"uploaded_at"`
	ProcessedAt time.Time `json:"-"`
}

type RequestWithdraw struct {
	Number  string `json:"order"`
	Accrual int    `json:"sum"`
}

type ResponseAccrual struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

// TODO можно было бы еще что-то типо такого сделать
//type OrderStatus string
//и потом
//var ( REGISTERED OrderStatus = "REGISTERED" ... )

const (
	StatusNew        = "NEW"        // Заказ загружен в систему, но не попал в обработку;
	StatusProcessing = "PROCESSING" // Вознаграждение за заказ рассчитывается;
	StatusInvalid    = "INVALID"    // Система расчёта вознаграждений отказала в расчёте;
	StatusProcessed  = "PROCESSED"  // Данные по заказу проверены и информация о расчёте успешно получена.
)
