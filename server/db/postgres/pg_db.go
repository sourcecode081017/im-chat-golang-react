/*
contains the implementation of the Postgres database connection and query execution
*/
package postgres

import (
	"context"
	"fmt"

	"github.com/sourcecode081017/im-chat-golang-react/models"
)

func (db *PgDb) CreateUser(ctx context.Context, user *models.User) error {
	// Insert a new user into the database
	query := `INSERT INTO users (username, first_name, last_name, email) VALUES ($1, $2, $3, $4) RETURNING id`
	var dbUser map[string]interface{}
	err := db.Conn.QueryRow(ctx, query, user.Username, user.FirstName, user.LastName, user.Email).Scan(dbUser["id"])
	// Check if the user already exists
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (db *PgDb) GetUsers(ctx context.Context) ([]models.User, error) {
	// Retrieve all users from the database
	query := `SELECT id, username, first_name, last_name, email FROM users`
	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
