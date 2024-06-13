package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework-1/internal/entities"
	"homework-1/internal/util"
	"log"
	"sync"
	"time"
)

// TODO: Remove Mutex, cause pgx.pool concurrency-safe

// SQLRepository field pool - concurrency-safe connection pool for pgx
type SQLRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
	Mu   sync.Mutex
}

func NewSQLRepository(ctx context.Context, cfg *entities.Config) *SQLRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			return err
		}
		return nil
	}, cfg.MaxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal(err)
	}
	return &SQLRepository{
		pool: pool,
		ctx:  ctx,
	}
}

func (r *SQLRepository) Insert(order entities.Order) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	query := `INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, hash) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(r.ctx, query, order.ID, order.UserID, order.StorageUntil, order.Issued, order.IssuedAt, order.Returned, order.Hash)
	return err
}

func (r *SQLRepository) Update(order entities.Order) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	query := `UPDATE orders SET user_id=$1, storage_until=$2, issued=$3, issued_at=$4, returned=$5, hash=$6
              WHERE id=$7`
	_, err := r.pool.Exec(r.ctx, query, order.UserID, order.StorageUntil, order.Issued, order.IssuedAt, order.Returned, order.Hash, order.ID)
	return err
}

func (r *SQLRepository) Exists(id string) (bool, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var exists bool
	err := r.pool.QueryRow(r.ctx, `SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)`, id).Scan(&exists)
	return exists, err
}

func (r *SQLRepository) Delete(id string) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	_, err := r.pool.Exec(context.Background(), `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLRepository) GetOrders() (map[string]entities.Order, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	rows, err := r.pool.Query(r.ctx, `SELECT id, user_id, storage_until, issued, issued_at, returned, hash FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[string]entities.Order)
	for rows.Next() {
		var order entities.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.StorageUntil, &order.Issued, &order.IssuedAt, &order.Returned, &order.Hash)
		if err != nil {
			return nil, err
		}
		orders[order.ID] = order
	}
	return orders, nil
}

func (r *SQLRepository) GetOrderIds() ([]string, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	query := `SELECT id FROM orders`
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderIDs []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		orderIDs = append(orderIDs, id)
	}
	return orderIDs, nil
}
