/*
Creates postgres connection pool and runs migration
*/

package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type PgDb struct {
	Conn *pgx.Conn
}

var DB_URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
	os.Getenv("PG_USER"),
	os.Getenv("PG_PASSWORD"),
	os.Getenv("PG_HOST"),
	os.Getenv("PG_PORT"),
	os.Getenv("PG_DB"),
)

var DB_QUERIES = []string{
	`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
		username VARCHAR(255) NOT NULL UNIQUE,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`,
	`CREATE TABLE IF NOT EXISTS channels (
		id SERIAL PRIMARY KEY,
		channel_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
		channel_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`,
	`CREATE TABLE IF NOT EXISTS user_channels (
		id SERIAL PRIMARY KEY,
		channel_id INT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (channel_id, user_id)
	);`,
}

func NewPgDb(ctx context.Context) (*PgDb, error) {

	conn, err := pgx.Connect(ctx, DB_URL)
	log.Printf("Connecting to database: %s\n", DB_URL)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Run migrations here if needed

	return &PgDb{Conn: conn}, nil
}

func (db *PgDb) Close(ctx context.Context) {
	if err := db.Conn.Close(ctx); err != nil {
		fmt.Printf("failed to close database connection: %v\n", err)
	}
}
func (db *PgDb) Ping(ctx context.Context) error {
	if err := db.Conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
func (db *PgDb) GetConn() *pgx.Conn {
	return db.Conn
}

func (db *PgDb) RunMigrations(ctx context.Context) error {
	for _, query := range DB_QUERIES {
		if _, err := db.Conn.Exec(context.Background(), query); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}
	return nil
}
