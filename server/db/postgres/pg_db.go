/*
contains the implementation of the Postgres database connection and query execution
*/
package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sourcecode081017/im-chat-golang-react/models"
	"gorm.io/gorm"
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

func (pg *PgDb) CreateChannel(ctx context.Context, channel *models.Channel) error {
	// Use a transaction to ensure atomicity
	return pg.Db.Transaction(func(tx *gorm.DB) error {
		// First, get the creator user by UUID
		var creator models.User
		if err := tx.Where("user_uuid = ?", channel.CreatedBy).First(&creator).Error; err != nil {
			return fmt.Errorf("failed to find creator user: %w", err)
		}

		// Create the channel
		if err := tx.Create(channel).Error; err != nil {
			return fmt.Errorf("failed to create channel: %w", err)
		}

		// Add the creator as a subscriber
		if err := tx.Model(channel).Association("Subscribers").Append(&creator); err != nil {
			return fmt.Errorf("failed to add creator as subscriber: %w", err)
		}

		return nil
	})
}

func (pg *PgDb) GetChannelByUUID(ctx context.Context, channelUUID uuid.UUID) (*models.Channel, error) {
	var channel models.Channel
	err := pg.Db.Preload("Subscribers").Where("channel_uuid = ?", channelUUID).First(&channel).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel: %w", err)
	}
	return &channel, nil
}

func (pg *PgDb) GetUserChannels(ctx context.Context, userUUID uuid.UUID) ([]models.Channel, error) {
	var user models.User
	err := pg.Db.Where("user_uuid = ?", userUUID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Get all channels where the user is a subscriber
	var channels []models.Channel
	err = pg.Db.Preload("Subscribers").
		Joins("JOIN channel_subscribers ON channel_subscribers.channel_id = channels.id").
		Where("channel_subscribers.user_id = ?", user.ID).
		Find(&channels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user's channels: %w", err)
	}

	return channels, nil
}
