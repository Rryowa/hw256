package db

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/internal/storage"
	"homework/internal/util"
	"strings"
)

type Repository struct {
	Pool      *pgxpool.Pool
	zapLogger *zap.SugaredLogger
}

func NewSQLRepository(ctx context.Context, cfg *config.DbConfig, zap *zap.SugaredLogger) storage.Storage {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			zap.Fatalln(err, "db connection error")
		}

		return nil
	}, cfg.Attempts, cfg.Timeout)

	if err != nil {
		zap.Fatalln(err, "DoWithTries error")
	}
	zap.Infoln("Connected to db")

	return &Repository{
		Pool:      pool,
		zapLogger: zap,
	}
}

func (r *Repository) Insert(ctx context.Context, order models.Order) (models.Order, error) {
	const op = "storage.Insert"
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Insert")
	defer span.Finish()

	query := `INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash`

	rows, err := r.Pool.Query(ctx, query, order.ID, order.UserID,
		order.StorageUntil, order.Issued, order.IssuedAt, order.Returned,
		order.OrderPrice, order.Weight, order.PackageType, order.PackagePrice,
		order.Hash)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var insertedOrder models.Order
	err = pgxscan.ScanOne(&insertedOrder, rows)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op2, err)
	}

	return insertedOrder, nil
}

func (r *Repository) Update(ctx context.Context, order models.Order) (models.Order, error) {
	const op = "storage.Update"
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Update")
	defer span.Finish()

	query := `UPDATE orders SET returned=$1
        WHERE id=$2 RETURNING id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
        `

	rows, err := r.Pool.Query(ctx, query, order.Returned, order.ID)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var updatedOrder models.Order
	err = pgxscan.ScanOne(&updatedOrder, rows)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op2, err)
	}

	return updatedOrder, nil
}

func (r *Repository) IssueUpdate(ctx context.Context, orders []models.Order) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.IssueUpdate")
	defer span.Finish()

	query := `UPDATE orders SET issued=$1, issued_at=NOW()
        WHERE id=$2
        `

	batch := &pgx.Batch{}
	for _, order := range orders {
		batch.Queue(query, order.Issued, order.ID)
		r.zapLogger.Infof("Order with id:%s issued\n", order.ID)
	}

	br := r.Pool.SendBatch(ctx, batch)
	for i := range orders {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return []models.Order{}, fmt.Errorf("error executing batch at order index %d: %w", i, err)
		}
	}
	err := br.Close()
	if err != nil {
		return []models.Order{}, err
	}

	var issuedOrders []models.Order
	for i, order := range orders {
		o, err := r.get(ctx, order.ID)
		if err != nil {
			return []models.Order{}, fmt.Errorf("error getting issued order by index %d: %w", i, err)
		}
		issuedOrders = append(issuedOrders, o)
	}

	return issuedOrders, err
}

func (r *Repository) Delete(ctx context.Context, id string) (string, error) {
	const op = "storage.Delete"
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Delete")
	defer span.Finish()

	query := `DELETE FROM orders WHERE id=$1 RETURNING id
		`

	var idd string
	err := r.Pool.QueryRow(ctx, query, id).Scan(&idd)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return idd, nil
}

func (r *Repository) GetReturns(ctx context.Context, offset, limit int) ([]models.Order, error) {
	const op = "storage.GetReturns"
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.GetReturns")
	defer span.Finish()

	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
        FROM orders
        WHERE returned = TRUE
        ORDER BY id
        OFFSET $1
 		FETCH NEXT $2 ROWS ONLY`

	rows, err := r.Pool.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var returns []models.Order
	if err := pgxscan.ScanAll(&returns, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	return returns, nil
}

func (r *Repository) GetOrders(ctx context.Context, userId string, offset, limit int) ([]models.Order, error) {
	const op = "storage.GetOrders"
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.GetOrders")
	defer span.Finish()

	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
		FROM orders
		WHERE user_id = $1 AND issued = FALSE
		ORDER BY storage_until
		OFFSET $2
		FETCH NEXT $3 ROWS ONLY`

	rows, err := r.Pool.Query(ctx, query, userId, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var userOrders []models.Order
	if err := pgxscan.ScanAll(&userOrders, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}
	return userOrders, err
}

func (r *Repository) InsertEvent(ctx context.Context, request string) (models.Event, error) {
	const op = "storage.InsertEvent"

	args := strings.Split(request, " ")
	rows, err := r.Pool.Query(ctx,
		`INSERT INTO events (method_name, request, status, requested_at) 
				VALUES ($1, $2, $3, NOW()) returning id, method_name, request, status, requested_at`,
		args[0], request, models.EventStatusRequested)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var event models.Event
	err = pgxscan.ScanOne(&event, rows)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op2, err)
	}

	return event, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, event models.Event) (models.Event, error) {
	const op = "UpdateEvent"

	rows, err := r.Pool.Query(ctx,
		`UPDATE events SET status = $1, processed_at = NOW()
				WHERE id = $2 returning id, method_name, request, status, requested_at, processed_at`,
		models.EventStatusProcessed, event.ID)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var updEvent models.Event
	err = pgxscan.ScanOne(&updEvent, rows)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op2, err)
	}

	return updEvent, nil
}

func (r *Repository) get(ctx context.Context, id string) (models.Order, error) {
	const op = "storage.Get"
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.get")
	defer span.Finish()

	var order models.Order
	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
				FROM orders
				WHERE id=$1`

	if err := pgxscan.Get(ctx, r.Pool, &order, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return models.Order{}, util.ErrOrderNotFound
		} else {
			return models.Order{}, fmt.Errorf("%s: %w", op, err)
		}
	}
	return order, nil
}

func (r *Repository) Truncate(ctx context.Context, table string) error {
	_, err := r.Pool.Exec(ctx, fmt.Sprintf(`truncate table %s`, table))
	if err != nil {
		return err
	}
	return nil
}