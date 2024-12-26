package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/okiww/billing-loan-system/internal/billing_config/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
)

type billingConfigRepository struct {
	*mysql.DBMySQL
}

// GetBillingConfigByName get a BillingConfig by name
func (repo *billingConfigRepository) GetBillingConfigByName(ctx context.Context, name string) (*models.BillingConfig, error) {

	billingConfig := &models.BillingConfig{}
	query := "SELECT name, value FROM billing_configs WHERE name = ?"
	err := repo.DB.GetContext(ctx, billingConfig, query, strings.ToLower(name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no config found") // No rows found
		}
		return nil, err
	}
	return billingConfig, nil
}

func NewBillingConfigRepository(db *mysql.DBMySQL) BillingConfigRepositoryInterface {
	return &billingConfigRepository{db}
}

type BillingConfigRepositoryInterface interface {
	GetBillingConfigByName(ctx context.Context, name string) (*models.BillingConfig, error)
}
