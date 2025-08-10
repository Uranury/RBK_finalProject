package db

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// InitDB connects to the database and returns a ready-to-use handle.
// It fails fast if the database isn't reachable.
func InitDB(driver, dsn string, migrationsPath string, logger *slog.Logger) (*sqlx.DB, error) {
	// Run migrations first with their own connection
	if err := RunMigrations(dsn, migrationsPath, logger); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create the main application database connection
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Set connection pool settings (recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return db, nil
}

// RunMigrations applies any pending migrations from the specified path.
// It creates its own database connection for migrations to avoid closing the main connection.
func RunMigrations(dsn string, migrationsPath string, logger *slog.Logger) error {
	// Create a separate connection for migrations
	migrationDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to create migration connection: %w", err)
	}
	defer migrationDB.Close() // This will only close the migration connection

	driver, err := postgres.WithInstance(migrationDB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationURL := migrationsPath
	if !strings.HasPrefix(migrationsPath, "file://") {
		migrationURL = "file://" + migrationsPath
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationURL,
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
