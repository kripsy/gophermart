package storage

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kripsy/gophermart/internal/gophermart/internal/config"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store interface {
	PutOrder(ctx context.Context, userName interface{}, number int64) (models.Order, error)
	GetOrders(ctx context.Context, username interface{}) ([]models.ResponseOrder, error)
	GetProcessingOrders(ctx context.Context) ([]models.ResponseOrder, error)
	GetNewOrders(ctx context.Context) ([]models.ResponseOrder, error)
	GetBalance(ctx context.Context, userName interface{}) (models.ResponseBalance, error)
	PutWithdraw(ctx context.Context, userName interface{}, number int64, accrual int) error
	GetWithdraws(ctx context.Context, userName interface{}) ([]models.ResponseWithdrawals, error)
	UpdateStatusOrder(ctx context.Context, number string, status string, accrual int) (models.ResponseOrder, error)
}

type DBStorage struct{}

var s Store = &DBStorage{}

func GetStorage() Store {
	return s
}

func (s *DBStorage) PutOrder(ctx context.Context, userName interface{}, number int64) (models.Order, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable to Close connect to database: %v\n", zap.String("msg", err.Error()))
		}
	}(conn, ctx)
	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return models.Order{}, err
	}

	var ID int64
	var Username string
	var Number int64
	var Status string
	var Accrual int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "INSERT INTO public.gophermart_order (username, number, status) VALUES ($1, $2, $3) ON CONFLICT (number) DO UPDATE SET number=EXCLUDED.number RETURNING *;", userName, number, models.StatusNew).Scan(&ID, &Username, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
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
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)
	if err != nil {
		l.Error("Unable to connect to database: ", zap.String("msg", err.Error()))
		return []models.ResponseOrder{}, err
	}

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where username=$1 and accrual >= 0 order by uploaded_at;", username)
	//lint:ignore SA5001 func (rows *baseRows) Close() {} does not return an error !
	//defer rows.Close() //nolint:all

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
		order.UploadedAt = UploadedAt.Time
		order.ProcessedAt = ProcessedAt.Time

		orders = append(orders, order)
	}

	rows.Close()

	return orders, nil
}

func (s *DBStorage) GetProcessingOrders(ctx context.Context) ([]models.ResponseOrder, error) {

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

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where status=$1 and accrual >= 0 order by uploaded_at;", models.StatusProcessing)
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
		order.Username = UserName
		order.Number = strconv.FormatInt(Number, 10)
		order.Status = Status
		order.Accrual = Accrual
		order.UploadedAt = UploadedAt.Time
		order.ProcessedAt = ProcessedAt.Time

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *DBStorage) GetNewOrders(ctx context.Context) ([]models.ResponseOrder, error) {

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

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where status=$1 and accrual >= 0 order by uploaded_at;", models.StatusNew)
	if err != nil {
		return []models.ResponseOrder{}, err
	}
	defer rows.Close()

	var orders []models.ResponseOrder

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
		order.UploadedAt = UploadedAt.Time
		order.ProcessedAt = ProcessedAt.Time

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *DBStorage) GetBalance(ctx context.Context, userName interface{}) (models.ResponseBalance, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)
	if err != nil {
		l.Error("Unable to connect to database: ", zap.String("msg", err.Error()))
		return models.ResponseBalance{}, err
	}

	var ID int64
	var Username string
	var Current int
	var Withdrawn int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "select * from public.gophermart_balance where username=$1;", userName).Scan(&ID, &Username, &Current, &Withdrawn, &UploadedAt, &ProcessedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ResponseBalance{}, models.ErrNoBalance()
	}

	if err != nil {
		return models.ResponseBalance{}, err
	}

	balance := models.ResponseBalance{}

	balance.ID = ID
	balance.Username = Username
	balance.Current = Current
	balance.Withdrawn = Withdrawn
	balance.UploadedAt = UploadedAt.Time
	balance.ProcessedAt = ProcessedAt.Time

	return balance, nil
}

func (s *DBStorage) PutWithdraw(ctx context.Context, userName interface{}, number int64, accrual int) error {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("failed to conn.Close(ctx)", zap.String("msg", err.Error()))
		}
	}(conn, ctx)

	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return err
	}
	tx, err := conn.Begin(ctx)
	defer func(tx pgx.Tx) {
		err := tx.Rollback(ctx)
		if err != nil {
			l.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
	}(tx)

	if err != nil {
		l.Error("failed to Begin Tx in PutWithdraw", zap.String("msg", err.Error()))
		return err
	}

	var ID int64
	var UserName string
	var Current int
	var Withdrawn int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = tx.QueryRow(ctx, "update gophermart_balance set current=current-$2, withdrawn=withdrawn+$2 where current>=$2 and username=$1 returning *;", userName, accrual).Scan(&ID, &UserName, &Current, &Withdrawn, &UploadedAt, &ProcessedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNoBalance()
	}

	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO public.gophermart_order (username, number, status, accrual) VALUES ($1, $2, $3, $4);", userName, number, models.StatusProcessed, -accrual)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { //duplicate key value violates unique constraint "gophermart_order_number_key"
				return models.ErrDuplicateOrder()
			}
		}

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return err
}

func (s *DBStorage) GetWithdraws(ctx context.Context, userName interface{}) ([]models.ResponseWithdrawals, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)
	if err != nil {
		l.Error("Unable to connect to database: ", zap.String("msg", err.Error()))
		return []models.ResponseWithdrawals{}, err
	}

	rows, err := conn.Query(ctx, "select * from public.gophermart_order where username=$1 and accrual < 0 order by uploaded_at;", userName)
	//lint:ignore SA5001 func (rows *baseRows) Close() {} does not return an error !
	//defer rows.Close() //nolint:all

	if errors.Is(err, pgx.ErrNoRows) {
		return []models.ResponseWithdrawals{}, models.ErrNoOrder()
	}

	if err != nil {
		return []models.ResponseWithdrawals{}, err
	}

	orders := []models.ResponseWithdrawals{}

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

		order := models.ResponseWithdrawals{}

		order.ID = ID
		order.Username = Username
		order.Number = strconv.FormatInt(Number, 10)
		order.Status = Status
		order.Accrual = -Accrual
		order.UploadedAt = UploadedAt.Time
		order.ProcessedAt = ProcessedAt.Time

		orders = append(orders, order)
	}

	rows.Close()

	return orders, nil
}

func (s *DBStorage) UpdateStatusOrder(ctx context.Context, number string, status string, accrual int) (models.ResponseOrder, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return models.ResponseOrder{}, err
	}
	defer conn.Close(ctx)

	tx, err := conn.Begin(ctx)
	defer func(tx pgx.Tx) {
		err := tx.Rollback(ctx)
		if err != nil {
			l.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
	}(tx)

	if err != nil {
		l.Error("failed to Begin Tx in PutWithdraw", zap.String("msg", err.Error()))
		return models.ResponseOrder{}, err
	}

	var ID int64
	var Username string
	var Number string
	var Status string
	var Accrual int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "update public.gophermart_order set status=$1, accrual=$2 where number=$3 returning *;", status, accrual, number).Scan(&ID, &Username, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			l.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
		return models.ResponseOrder{}, err
	}
	_, err = tx.Exec(ctx, "insert into gophermart_balance (username, current, withdrawn) values ($1, $2, $3) on conflict (username) do update set current=gophermart_balance.current+$2;", Username, accrual, 0)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			l.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
		return models.ResponseOrder{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			l.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
		return models.ResponseOrder{}, err
	}

	order := models.ResponseOrder{}

	order.ID = ID
	order.Username = Username
	order.Number = Number
	order.Status = Status
	order.Accrual = Accrual
	order.UploadedAt = UploadedAt.Time
	order.ProcessedAt = ProcessedAt.Time

	return order, nil
}
