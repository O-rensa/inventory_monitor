package initializers

import (
	"database/sql"
	"os"
)

func InitializePostGres() (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", os.Getenv("POSTGRES_DSN"))
	return sqlDB, err
}
