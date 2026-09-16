package database

import (
	"database/sql"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func LoadDB(driver string, userdb string, password string, dbname string, logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open(driver, userdb+":"+password+"@/"+dbname)
	if err != nil {
		return nil, err
	}
	logger.Info("database loaded", "driver", driver, "dbname", dbname)
	return db, nil
}

func MigrateDB(db *sql.DB, logger *slog.Logger) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("database migrated")

	return nil
}
