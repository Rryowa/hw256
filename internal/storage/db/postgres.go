package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"homework-1/internal/models"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"log"
	"strings"
	"time"
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
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			log.Fatal(err, "db connection error")
		}

		return nil
	}, cfg.MaxAttempts, 5*time.Second)

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

func (r *repository) ListReturns(limit, offset int) ([]models.Order, error) {
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

func (r *repository) ListOrders(userId string, limit int) ([]models.Order, error) {
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

func (r *repository) ApplyMigrations(direction string) error {
	db := stdlib.OpenDBFromPool(r.pool)
	defer db.Close()

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Unable to create database instance: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./schema",
		"cli", instance)
	if err != nil {
		log.Fatalf("Unable to create migrate instance: %v\n", err)
	}

	switch direction {
	case "up":
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Unable to apply up migration: %v\n", err)
		}
		log.Println("Up migration applied successfully!")
	case "down":
		err = m.Down()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Unable to apply down migration: %v\n", err)
		}
		log.Println("Down migration applied successfully!")
	default:
		log.Fatalf("Invalid migration direction: %s\n", direction)
	}

	return nil
}

func (r *repository) AnalyzeQueryPlan(query string, args ...interface{}) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	explainQuery := "EXPLAIN (ANALYZE, VERBOSE) " + query
	rows, err := conn.Query(r.ctx, explainQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	s := strings.ReplaceAll(query, "\t", "")
	ss := strings.Split(s, " ")
	fmt.Println(ss)
	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return err
		}
		fmt.Println(plan)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
