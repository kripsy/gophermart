package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kripsy/gophermart/internal/gophermart/internal/config"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct{}

var s = DBStorage{}

func GetStorage() DBStorage {
	return s
}

func (s *DBStorage) PutOrder(ctx context.Context, userName interface{}, number int64) (models.Order, error) {

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
	var Username string
	var Number int64
	var Status string
	var Accrual int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "INSERT INTO public.gophermart_order (username, number, status) VALUES ($1, $2, $3) ON CONFLICT (number) DO UPDATE SET number=EXCLUDED.number RETURNING gophermart_order.id, gophermart_order.username, gophermart_order.number, gophermart_order.status, gophermart_order.accrual, gophermart_order.uploaded_at, gophermart_order.processed_at;", userName, number, models.StatusNew).Scan(&ID, &Username, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{}

	order.ID = ID
	order.Username = Username
	order.Number = Number
	order.Status = Status
	order.Accrual = Accrual
	order.UploadedAt = UploadedAt
	order.ProcessedAt = ProcessedAt

	return order, nil
}

func (s *DBStorage) GetOrders(ctx context.Context, username interface{}) ([]models.ResponseOrder, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	if err != nil {
		l.Error("Unable to connect to database: ", zap.String("msg", err.Error()))
		return []models.ResponseOrder{}, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where username=$1 and accrual >= 0 order by uploaded_at;", username)
	//lint:ignore SA5001 should check returned wraperror before deferring rows.Close()
	defer rows.Close()

	if err != nil {
		return []models.ResponseOrder{}, err
	}

	orders := make([]models.ResponseOrder, 0)

	for rows.Next() {
		var ID int64
		var Username string
		var Number int64
		var Status string
		var Accrual int
		var UploadedAt pgtype.Timestamptz
		var ProcessedAt pgtype.Timestamptz
		err = rows.Scan(&ID, &Username, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
		if err != nil {
			return nil, err
		}

		order := models.ResponseOrder{}

		order.ID = ID
		order.Username = Username
		order.Number = strconv.FormatInt(Number, 10)
		order.Status = Status
		order.Accrual = Accrual
		order.UploadedAt = UploadedAt
		order.ProcessedAt = ProcessedAt

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *DBStorage) GetBalance(ctx context.Context, userName interface{}) (models.ResponseBalance, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	if err != nil {
		l.Error("Unable to connect to database: ", zap.String("msg", err.Error()))
		return models.ResponseBalance{}, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)

	var ID int64
	var Username string
	var Current int
	var Withdrawn int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "select * from public.gophermart_balance where username=$1;", userName).Scan(&ID, &Username, &Current, &Withdrawn, &UploadedAt, &ProcessedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ResponseBalance{}, fmt.Errorf("%s %w", models.ErrUserOrdersNotRegistered, err)
	}

	if err != nil {
		return models.ResponseBalance{}, err
	}

	balance := models.ResponseBalance{}

	balance.ID = ID
	balance.Username = Username
	balance.Current = Current
	balance.Withdrawn = Withdrawn
	balance.UploadedAt = UploadedAt
	balance.ProcessedAt = ProcessedAt

	return balance, nil
}
