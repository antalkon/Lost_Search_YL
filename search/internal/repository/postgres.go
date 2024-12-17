package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/postgres"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(cfg postgres.Config) (*PostgresRepo, error) {
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresRepo: %w", err)
	}
	return &PostgresRepo{
		db: db,
	}, nil
}

func (r *PostgresRepo) CreateFind() {

}

func (r *PostgresRepo) GetFind() {

}
