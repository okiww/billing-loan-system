package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/okiww/billing-loan-system/helpers"
	"github.com/okiww/billing-loan-system/internal/user/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"sync"
)

var (
	repo     UserRepositoryInterface
	repoLock sync.Once
)

type userRepository struct {
	*mysql.DBMySQL
}

func NewUserRepository(db *mysql.DBMySQL) UserRepositoryInterface {
	if helpers.IsTestEnv() { // Skip singleton in tests
		return &userRepository{
			db,
		}
	}

	repoLock.Do(func() {
		repo = &userRepository{
			db,
		}
	})
	return repo
}

// UpdateUserToDelinquent updates the is_delinquent field of a user to true.
func (u *userRepository) UpdateUserToDelinquent(ctx context.Context, userID int32) error {
	query := `
		UPDATE user
		SET is_delinquent = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := u.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error updating user to delinquent: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by their ID.
func (u *userRepository) GetUserByID(ctx context.Context, userID int32) (*models.UserModel, error) {
	query := `
		SELECT id, name, is_delinquent
		FROM user
		WHERE id = ?
	`

	user := &models.UserModel{}
	err := u.DB.GetContext(ctx, user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with the given ID
		}
		return nil, fmt.Errorf("error retrieving user by ID: %w", err)
	}

	return user, nil
}

type UserRepositoryInterface interface {
	UpdateUserToDelinquent(ctx context.Context, userID int32) error
	GetUserByID(ctx context.Context, userID int32) (*models.UserModel, error)
}
