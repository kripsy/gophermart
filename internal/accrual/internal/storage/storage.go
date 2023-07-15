package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kripsy/gophermart/internal/accrual/internal/config"
	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	"github.com/kripsy/gophermart/internal/accrual/internal/models"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct{}

var s = DBStorage{}

func GetStorage() DBStorage {
	return s
}

func (s *DBStorage) PutOrder(ctx context.Context, number int64) (models.Order, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return models.Order{}, err
	}
	defer conn.Close(ctx)

	var ID int64
	var Number int64
	var Status string
	var Accural int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "INSERT INTO public.accrual (number, status, accrual) VALUES ($1, $2, $3) ON CONFLICT (number) DO UPDATE SET number=EXCLUDED.number RETURNING accrual.id, accrual.number, accrual.status, accrual.accrual, accrual.uploaded_at, accrual.processed_at;", number, models.StatusProcessed, number%1000).Scan(&ID, &Number, &Status, &Accural, &UploadedAt, &ProcessedAt)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{}

	order.ID = ID
	order.Number = Number
	order.Status = Status
	order.Accural = Accural
	order.UploadedAt = UploadedAt
	order.ProcessedAt = ProcessedAt

	return order, nil
}
