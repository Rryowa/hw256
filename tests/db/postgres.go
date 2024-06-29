package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/models"
	"homework/internal/storage/db"
	"log"
)

type TestRepository struct {
	Repo *db.Repository
}

func NewTestRepository(ctx context.Context, cfg *models.Config, schemaName string) TestRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err, "db connection error")
	}
	log.Println("Connected to db")

	return TestRepository{
		&db.Repository{
			Pool: pool,
			Ctx:  ctx,
		},
	}
}

//	func loadTestData(t *testing.T, db *sql.DB, schemaName string, testDatabase string) {
//		for _, testDataName := range testDataNames {
//			file, err := os.Open(fmt.Sprintf("./testdata/%s.sql", testDataName))
//			require.NoError(t, err)
//			reader := bufio.NewReader(file)
//			var query string
//			for {
//				line, err := reader.ReadString('\n')
//				if err == io.EOF {
//					break
//				}
//				require.NoError(t, err)
//				line = line[:len(line)-1]
//				if line == "" {
//					query = addSchemaPrefix(schemaName, query)
//					_, err := db.Exec(query)
//					require.NoError(t, err)
//					query = ""
//				}
//				query += line
//			}
//			file.Close()
//		}
//	}
