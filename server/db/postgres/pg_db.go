/*
contains the implementation of the Postgres database connection and query execution
*/
package postgres

import (
	"context"
	"fmt"

	"github.com/sourcecode081017/im-chat-golang-react/models"
)

func (pg *PgDb) CreateUser(ctx context.Context, user *models.User) error {
	// Insert a new user into the database
	err := pg.Db.Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (pg *PgDb) GetUsers(ctx context.Context) ([]models.User, error) {
	// Retrieve all users from the database
	var users []models.User
	err := pg.Db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

func (pg *PgDb) CreateChannel(ctx context.Context, channel *models.Channel) error {
	// Insert a new channel into the database
	err := pg.Db.Create(channel).Error
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	return nil
}
