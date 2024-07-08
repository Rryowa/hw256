package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"log"
	"strings"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func NewSQLRepository(ctx context.Context, cfg *models.Config) storage.Storage {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			log.Fatal(err, "db connection error")
		}

		return nil
	}, cfg.Attempts, cfg.Timeout)

	if err != nil {
		log.Fatal(err, "DoWithTries error")
	}
	log.Println("Connected to db")

	return &Repository{
		Pool: pool,
	}
}

func (r *Repository) Insert(ctx context.Context, order models.Order) (string, error) {
	query := `INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`

	var id string
	row := r.Pool.QueryRow(ctx, query, order.ID, order.UserID,
		order.StorageUntil, order.Issued, order.IssuedAt, order.Returned,
		order.OrderPrice, order.Weight, order.PackageType, order.PackagePrice,
		order.Hash)
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return "", err
	}
	return id, nil
}

func (r *Repository) Update(ctx context.Context, order models.Order) (bool, error) {
	query := `UPDATE orders SET returned=$1
        WHERE id=$2 RETURNING returned
        `

	var returned bool
	err := r.Pool.QueryRow(ctx, query, order.Returned, order.ID).Scan(&returned)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return false, err
	}
	return returned, nil
}

func (r *Repository) IssueUpdate(ctx context.Context, orders []models.Order) ([]bool, error) {
	query := `UPDATE orders SET issued=$1, issued_at=NOW()
        WHERE id=$2
        `

	batch := &pgx.Batch{}
	for _, order := range orders {
		batch.Queue(query, order.Issued, order.ID)
		log.Printf("Order with id:%s issued\n", order.ID)
	}

	br := r.Pool.SendBatch(ctx, batch)
	var issuedOrders []bool
	for i, order := range orders {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return []bool{}, fmt.Errorf("error executing batch at order index %d: %w", i, err)
		}
		issuedOrders = append(issuedOrders, order.Issued)
	}
	err := br.Close()
	if err != nil {
		return []bool{}, err
	}
	return issuedOrders, err
}

func (r *Repository) Delete(ctx context.Context, id string) (string, error) {
	query := `DELETE FROM orders WHERE id=$1 RETURNING id
		`

	var idd string
	err := r.Pool.QueryRow(ctx, query, id).Scan(&idd)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return "", err
	}

	return idd, nil
}

func (r *Repository) Get(ctx context.Context, id string) (models.Order, error) {
	var order models.Order
	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash FROM orders
		WHERE id=$1
		`

	if err := pgxscan.Get(ctx, r.Pool, &order, query, id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
			return models.Order{}, err
		} else if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, util.ErrOrderNotFound
		}
	}
	return order, nil
}

func (r *Repository) GetReturns(ctx context.Context, offset, limit int) ([]models.Order, error) {
	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
        FROM orders
        WHERE returned = TRUE
        ORDER BY id
        OFFSET $1
 		FETCH NEXT $2 ROWS ONLY
    `

	rows, err := r.Pool.Query(ctx, query, offset, limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}

	defer rows.Close()

	var returns []models.Order
	if err := pgxscan.ScanAll(&returns, rows); err != nil {
		return nil, err
	}

	return returns, nil
}

func (r *Repository) GetOrders(ctx context.Context, userId string, offset, limit int) ([]models.Order, error) {
	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
		FROM orders
		WHERE user_id = $1 AND issued = FALSE
		ORDER BY storage_until
		OFFSET $2
		FETCH NEXT $3 ROWS ONLY
	`

	rows, err := r.Pool.Query(ctx, query, userId, offset, limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	defer rows.Close()

	var userOrders []models.Order
	if err := pgxscan.ScanAll(&userOrders, rows); err != nil {
		return nil, err
	}
	return userOrders, err
}

func (r *Repository) Truncate(ctx context.Context, table string) error {
	_, err := r.Pool.Exec(ctx, fmt.Sprintf(`truncate table %s`, table))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateEvent(ctx context.Context, request string) error {
	const op = "CreateEvent"
	args := strings.Split(request, " ")
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO events (method_name, request, status, acquired_at) 
				VALUES ($1, $2, $3, NOW())`,
		args[0], request, models.EventStatusAcquired)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) GetEvent(ctx context.Context) (models.Event, pgx.Tx, error) {
	const op = "GetEvent"
	tx, _ := r.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	row := tx.QueryRow(ctx,
		`SELECT id, request, method_name, status, acquired_at
				FROM events WHERE status = $1`, models.EventStatusAcquired)
	var e models.Event
	err := row.Scan(&e.ID, &e.Request, &e.MethodName, &e.Status, &e.AcquiredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e, nil, nil
		}
		return e, nil, fmt.Errorf("%s: %w", op, err)
	}

	return e, tx, nil
}

func (r *Repository) ProcessEvent(ctx context.Context, e models.Event, tx pgx.Tx) (models.Event, error) {
	const op = "ProcessEvent"
	e.Status = models.EventStatusProcessed
	row := tx.QueryRow(ctx,
		`UPDATE events SET status = $1, processed_at = NOW()
				WHERE id = $2 returning processed_at`,
		e.Status, e.ID)
	err := row.Scan(&e.ProcessedAt)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Event{}, err
	}

	return e, nil
}