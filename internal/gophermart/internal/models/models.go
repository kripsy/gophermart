package models

import "github.com/jackc/pgx/v5/pgtype"

type Order struct {
	ID          int64
	UserName    string
	Number      int64
	Status      string
	Accrual     int
	UploadedAt  pgtype.Timestamptz
	ProcessedAt pgtype.Timestamptz
}

type ResponseOrder struct {
	ID          int64              `json:"-"`
	UserName    string             `json:"-"`
	Number      string             `json:"order"`
	Status      string             `json:"status"`
	Accrual     int                `json:"accrual,omitempty"`
	UploadedAt  pgtype.Timestamptz `json:"uploaded_at"`
	ProcessedAt pgtype.Timestamptz `json:"-"`
}

type ResponseBalance struct {
	ID          int64              `json:"-"`
	UserName    string             `json:"-"`
	Current     int                `json:"current"`
	Withdrawn   int                `json:"withdrawn"`
	UploadedAt  pgtype.Timestamptz `json:"uploaded_at"`
	ProcessedAt pgtype.Timestamptz `json:"-"`
}

type RequestWithdraw struct {
	Number  int64 `json:"order"`
	Accrual int   `json:"sum"`
}

const (
	StatusNew        = "NEW"        // Заказ загружен в систему, но не попал в обработку;
	StatusProcessing = "PROCESSING" // Вознаграждение за заказ рассчитывается;
	StatusInvalid    = "INVALID"    // Система расчёта вознаграждений отказала в расчёте;
	StatusProcessed  = "PROCESSED"  // Данные по заказу проверены и информация о расчёте успешно получена.
)
