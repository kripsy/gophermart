package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kripsy/gophermart/internal/gophermart/internal/config"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct{}

func (s DBStorage) InitStorage() error {
	//TODO implement me
	panic("implement me")
}

var s = DBStorage{}

func GetStorage() DBStorage {
	return s
}

func (s *DBStorage) PutOrder(ctx context.Context, userName interface{}, number int) (models.Order, error) {

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
	var UserName string
	var Number int64
	var Status string
	var Accural int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "INSERT INTO public.gophermart_order (username, number, status) VALUES ($1, $2, $3) ON CONFLICT (number) DO UPDATE SET number=EXCLUDED.number RETURNING gophermart_order.id, gophermart_order.username, gophermart_order.number, gophermart_order.status, gophermart_order.accrual, gophermart_order.uploaded_at, gophermart_order.processed_at;", userName, number, models.StatusNew).Scan(&ID, &UserName, &Number, &Status, &Accural, &UploadedAt, &ProcessedAt)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{}

	order.ID = ID
	order.UserName = UserName
	order.Number = Number
	order.Status = Status
	order.Accural = Accural
	order.UploadedAt = UploadedAt
	order.ProcessedAt = ProcessedAt

	return order, nil
}
