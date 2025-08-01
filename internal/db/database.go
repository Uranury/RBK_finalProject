package db

import "database/sql"

func InitDB(driver string, dsn string) (*sqlx.DB, error) {}
