package storage

import (
	"context"
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
	var UserName string
	var Number int64
	var Status string
	var Accrual int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "INSERT INTO public.gophermart_order (username, number, status) VALUES ($1, $2, $3) ON CONFLICT (number) DO UPDATE SET number=EXCLUDED.number RETURNING *;", userName, number, models.StatusNew).Scan(&ID, &UserName, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{}

	order.ID = ID
	order.UserName = UserName
	order.Number = Number
	order.Status = Status
	order.Accrual = Accrual
	order.UploadedAt = UploadedAt
	order.ProcessedAt = ProcessedAt

	return order, nil
}

func (s *DBStorage) GetOrders(ctx context.Context, userName interface{}) ([]models.ResponseOrder, error) {

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

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where username=$1 and accrual >= 0 order by uploaded_at;", userName)
	if err != nil {
		return []models.ResponseOrder{}, err
	}
	defer rows.Close()

	var orders []models.ResponseOrder

	for rows.Next() {
		var ID int64
		var UserName string
		var Number int64
		var Status string
		var Accrual int
		var UploadedAt pgtype.Timestamptz
		var ProcessedAt pgtype.Timestamptz
		err = rows.Scan(&ID, &UserName, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
		if err != nil {
			return nil, err
		}

		order := models.ResponseOrder{}

		order.ID = ID
		order.UserName = UserName
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
	var UserName string
	var Current int
	var Withdrawn int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "select * from public.gophermart_balance where username=$1;", userName).Scan(&ID, &UserName, &Current, &Withdrawn, &UploadedAt, &ProcessedAt)
	if err != nil {
		return models.ResponseBalance{}, err
	}

	balance := models.ResponseBalance{}

	balance.ID = ID
	balance.UserName = UserName
	balance.Current = Current
	balance.Withdrawn = Withdrawn
	balance.UploadedAt = UploadedAt
	balance.ProcessedAt = ProcessedAt

	return balance, nil
}

func (s *DBStorage) PutWithdraw(ctx context.Context, userName interface{}, number int64, accrual int) error {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("failed to conn.Close(ctx)", zap.String("msg", err.Error()))
		}
	}(conn, ctx)

	tx, err := conn.Begin(ctx)
	if err != nil {
		l.Error("failed to Begin Tx in PutWithdraw", zap.String("msg", err.Error()))
		return err
	}

	defer func(tx pgx.Tx) {
		err := tx.Rollback(ctx)
		if err != nil {
			l.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
	}(tx)

	var ID int64
	var UserName string
	var Current int
	var Withdrawn int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = tx.QueryRow(ctx, "update gophermart_balance set current=current-$2, withdrawn=withdrawn+$2 where current>$2 and username=$1 returning *;", userName, accrual).Scan(&ID, &UserName, &Current, &Withdrawn, &UploadedAt, &ProcessedAt)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO public.gophermart_order (username, number, status, accrual) VALUES ($1, $2, $3, $4);", userName, number, models.StatusProcessed, -accrual)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return err
}

func (s *DBStorage) GetWithdraws(ctx context.Context, userName interface{}) ([]models.ResponseOrder, error) {

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

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where username=$1 and accrual < 0 order by uploaded_at;", userName)
	if err != nil {
		return []models.ResponseOrder{}, err
	}
	defer rows.Close()

	var orders []models.ResponseOrder

	for rows.Next() {
		var ID int64
		var UserName string
		var Number int64
		var Status string
		var Accrual int
		var UploadedAt pgtype.Timestamptz
		var ProcessedAt pgtype.Timestamptz
		err = rows.Scan(&ID, &UserName, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
		if err != nil {
			return nil, err
		}

		order := models.ResponseOrder{}

		order.ID = ID
		order.UserName = UserName
		order.Number = strconv.FormatInt(Number, 10)
		order.Status = Status
		order.Accrual = -Accrual
		order.UploadedAt = UploadedAt
		order.ProcessedAt = ProcessedAt

		orders = append(orders, order)
	}

	return orders, nil
}
