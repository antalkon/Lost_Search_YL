package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	UserName string
	Password string
	Host     string
	Port     string
	DbName   string
}

func New(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable hostname=%s port=%s", cfg.UserName, cfg.Password, cfg.DbName, cfg.Host, cfg.Port)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection: %w", err)
	}
	if _, err := db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return db, nil
}
