package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

// InitDB connects to the database and returns a ready-to-use handle.
// It fails fast if the database isn't reachable.
func InitDB(driver, dsn string, migrationsPath string, logger *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	if err := RunMigrations(db, migrationsPath, logger); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Database connection established successfully")
	return db, nil
}

// RunMigrations applies any pending migrations from the specified path.
// It returns nil if the database is already up to date.
func RunMigrations(db *sqlx.DB, migrationsPath string, logger *slog.Logger) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logger.Error("migration source close error", "error", sourceErr)
		}
		if dbErr != nil {
			logger.Error("migration database close error", "error", dbErr)
		}
	}()

	upErr := m.Up()
	if upErr != nil && !errors.Is(upErr, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", upErr)
	}

	if errors.Is(upErr, migrate.ErrNoChange) {
		logger.Info("Database is already up to date")
	} else {
		logger.Info("Database migrations completed successfully")
	}

	return nil
}
