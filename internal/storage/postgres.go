package storage

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework-1/internal/entities"
	"homework-1/internal/util"
	"log"
	"time"
)

// SQLRepository field pool - concurrency-safe connection pool for pgx
type SQLRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewSQLRepository(ctx context.Context, cfg *entities.Config) *SQLRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctxTimeout, connStr)
		if err != nil {
			log.Fatal(err, "db connection error")
		}
		return nil
	}, cfg.MaxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal(err, "DoWithTries error")
	}
	log.Println("Connected to db")
	return &SQLRepository{
		pool: pool,
		ctx:  ctx,
	}
}

func (r *SQLRepository) Insert(order entities.Order) error {
	query := `INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, hash) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(r.ctx, query, order.ID, order.UserID, order.StorageUntil, order.Issued, order.IssuedAt, order.Returned, order.Hash)

	return err
}

func (r *SQLRepository) Update(order entities.Order) error {
	query := `UPDATE orders SET issued=$1, issued_at=$2, returned=$3
              WHERE id=$4`
	_, err := r.pool.Exec(r.ctx, query, order.Issued, order.IssuedAt, order.Returned, order.ID)
	return err
}

func (r *SQLRepository) IssueUpdate(orders []entities.Order) error {
	tx, err := r.pool.Begin(r.ctx)
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
	err = br.Close()
	if err != nil {
		return err
	}
	return tx.Commit(r.ctx)
}

func (r *SQLRepository) Delete(id string) error {
	_, err := r.pool.Exec(r.ctx, `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLRepository) ListReturns(limit, offset int) ([]entities.Order, error) {
	query := `
        SELECT id, user_id, storage_until, issued, issued_at, returned
        FROM orders
        WHERE returned = TRUE
        ORDER BY id
        LIMIT $1 OFFSET $2
    `
	rows, err := r.pool.Query(r.ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var returns []entities.Order
	if err := pgxscan.ScanAll(&returns, rows); err != nil {
		return nil, err
	}
	return returns, nil
}

func (r *SQLRepository) ListOrders(userId string, limit int) ([]entities.Order, error) {
	query := `
        SELECT id, user_id, issued, storage_until, returned
        FROM orders
        WHERE user_id = $1 AND issued = FALSE
        ORDER BY storage_until DESC
        LIMIT $2
    `

	rows, err := r.pool.Query(r.ctx, query, userId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userOrders []entities.Order
	if err := pgxscan.ScanAll(&userOrders, rows); err != nil {
		return nil, err
	}

	return userOrders, nil
}
