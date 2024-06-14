package storage

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
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
	query := `UPDATE orders SET storage_until=$1, issued=$2, issued_at=$3, returned=$4
              WHERE id=$5`
	_, err := r.pool.Exec(r.ctx, query, order.StorageUntil, order.Issued, order.IssuedAt, order.Returned, order.ID)
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

	_, err := r.pool.Exec(r.ctx, `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLRepository) Get(id string) (entities.Order, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var order entities.Order
	query := `SELECT id, user_id, storage_until, issued, issued_at, returned, hash FROM orders WHERE id=$1`
	//err := r.pool.QueryRow(r.ctx, query, id).Scan(&order.ID, &order.UserID, &order.StorageUntil, &order.Issued, &order.IssuedAt, &order.Returned, &order.Hash)
	//if err != nil {
	//	return entities.Order{}, err
	//}
	if err := pgxscan.Get(r.ctx, r.pool, &order, query, id); err != nil {
		return entities.Order{}, err
	}
	return order, nil
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
