/*
contains the implementation of the Postgres database connection and query execution
*/
package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (pg *PgDb) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*models.User, error) {
	var user models.User
	err := pg.Db.Where("user_uuid = ?", userUUID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by UUID: %w", err)
	}
	return &user, nil
}

func (pg *PgDb) SaveMessage(ctx context.Context, message *models.ChatMessage) error {
	err := pg.Db.Create(message).Error
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

func (pg *PgDb) GetMessagesBetweenUsers(ctx context.Context, userId1, userId2 uuid.UUID) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := pg.Db.Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		userId1, userId2, userId2, userId1,
	).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	return messages, nil
}
