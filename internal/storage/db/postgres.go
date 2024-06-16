package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework-1/internal/models"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"log"
)

// repository field pool - concurrency-safe connection pool for pgx
type repository struct {
	pool *pgxpool.Pool
	ctx  context.Context
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

	return &repository{
		pool: pool,
		ctx:  ctx,
	}
}

func (r *repository) Insert(order models.Order) error {
	query := `INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, hash) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(r.ctx, query, order.ID, order.UserID, order.StorageUntil, order.Issued, order.IssuedAt, order.Returned, order.Hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return err
	}
	return nil
}

func (r *repository) Update(order models.Order) error {
	query := `UPDATE orders SET issued=$1, issued_at=$2, returned=$3
              WHERE id=$4`

	_, err := r.pool.Exec(r.ctx, query, order.Issued, order.IssuedAt, order.Returned, order.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return err
	}
	return nil
}

func (r *repository) IssueUpdate(orders []models.Order) error {
	tx, err := r.pool.BeginTx(r.ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(r.ctx)

	query := `UPDATE orders SET issued=$1, issued_at=$2, returned=$3
              WHERE id=$4`
	batch := &pgx.Batch{}
	for _, order := range orders {
		batch.Queue(query, order.Issued, order.IssuedAt, order.Returned, order.ID)
		log.Printf("Order with id:%s issued\n", order.ID)
	}

	br := tx.SendBatch(r.ctx, batch)
	for _, i := range orders {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("error executing batch at order index %d: %w", i, err)
		}
	}
	err = br.Close()

	return tx.Commit(r.ctx)
}

func (r *repository) Delete(id string) error {
	query := `DELETE FROM orders WHERE id=$1`

	_, err := r.pool.Exec(r.ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return err
	}

	return nil
}

func (r *repository) Exists(id string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)`
	if err := pgxscan.Get(r.ctx, r.pool, &exists, query, id); err != nil {
		log.Println(err)
	}
	return exists
}

func (r *repository) Get(id string) models.Order {
	var order models.Order
	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, hash FROM orders WHERE id=$1`
	if err := pgxscan.Get(r.ctx, r.pool, &order, query, id); err != nil {
		log.Println(err)
		return models.Order{}
	}
	return order
}

func (r *repository) GetReturns(limit, offset int) ([]models.Order, error) {
	query := `
        SELECT id, user_id, storage_until, issued, issued_at, returned
        FROM orders
        WHERE returned = TRUE
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	rows, err := r.pool.Query(r.ctx, query, limit, offset)
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

func (r *repository) GetOrders(userId string, limit int) ([]models.Order, error) {
	query := `
        SELECT id, user_id, issued, storage_until, returned
        FROM orders
        WHERE user_id = $1 AND issued = FALSE
        ORDER BY storage_until DESC
        LIMIT $2
    `

	rows, err := r.pool.Query(r.ctx, query, userId, limit)
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
