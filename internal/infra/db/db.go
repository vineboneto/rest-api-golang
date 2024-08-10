package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Conn *pgxpool.Pool
}

func BuildDSN() string {

	port, err := strconv.ParseUint(os.Getenv("POSTGRES_PORT"), 10, 16)

	if err != nil {
		panic("POSTGRES_PORT is not int")

	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASS"),
		os.Getenv("POSTGRES_HOST"),
		uint16(port),
		os.Getenv("POSTGRES_DATABASE"),
	)

	return dsn
}

func NewPostgresDB() (*DB, error) {

	config, err := pgxpool.ParseConfig(BuildDSN())
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	config.MaxConns = 15                      // Número máximo de conexões no pool
	config.MaxConnLifetime = 30 * time.Minute // Tempo máximo de vida de uma conexão
	config.MaxConnIdleTime = 5 * time.Minute  // Tempo máximo de inatividade de uma conexão
	config.MinConns = 2

	conn, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	return &DB{Conn: conn}, nil
}

func (db *DB) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *DB) Close() {
	db.Conn.Close()
}
