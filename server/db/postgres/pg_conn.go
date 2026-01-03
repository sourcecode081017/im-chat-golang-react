/*
Creates postgres connection pool and runs migration
*/

package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/sourcecode081017/im-chat-golang-react/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDb struct {
	Db *gorm.DB
}

var DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("PG_HOST"),
	os.Getenv("PG_PORT"),
	os.Getenv("PG_USER"),
	os.Getenv("PG_PASSWORD"),
	os.Getenv("PG_DB"),
)

func NewPgDb(ctx context.Context) (*PgDb, error) {

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &PgDb{Db: db}, nil
}

func (pg *PgDb) RunMigrations(ctx context.Context) error {
	if err := pg.Db.AutoMigrate(&models.User{}, &models.ChatMessage{}); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	err := pg.Db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_email ON users (email) WHERE email IS NOT NULL;").Error
	if err != nil {
		return fmt.Errorf("failed to create unique index on email: %w", err)
	}
	return nil
}
