package main

import (
	"context"
	"homework-1/internal/storage/db"
	"homework-1/internal/util"
	"log"
	"time"
)

func main() {
	repo := db.NewSQLRepository(context.Background(), util.NewConfig())
	if err := repo.ApplyMigrations("up"); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}
	// Анализ Insert
	err := repo.AnalyzeQueryPlan(
		`INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, hash) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"some-id", "some-user-id", time.Now(), false, time.Now(), false, "some-hash",
	)

	if err != nil {
		log.Fatal(err)
	}

	// Анализ Update
	err = repo.AnalyzeQueryPlan(
		`UPDATE orders SET issued=$1, issued_at=$2, returned=$3 WHERE id=$4`,
		true, time.Now(), false, "some-id",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Анализ IssueUpdate
	err = repo.AnalyzeQueryPlan(
		`UPDATE orders SET issued=$1, issued_at=$2, returned=$3 WHERE id=$4`,
		true, time.Now(), false, "some-id",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Анализ Delete
	err = repo.AnalyzeQueryPlan(
		`DELETE FROM orders WHERE id=$1`,
		"some-id",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Анализ ListReturns
	err = repo.AnalyzeQueryPlan(
		`SELECT id, user_id, storage_until, issued, issued_at, returned
			FROM orders
			WHERE returned = TRUE
			ORDER BY id
			LIMIT $1 OFFSET $2`,
		10, 0,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Анализ ListOrders
	err = repo.AnalyzeQueryPlan(
		`SELECT id, user_id, issued, storage_until, returned
			FROM orders
			WHERE user_id = $1 AND issued = FALSE
			ORDER BY storage_until DESC
			LIMIT $2`,
		"some-user-id", 10,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := repo.ApplyMigrations("down"); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}
}
