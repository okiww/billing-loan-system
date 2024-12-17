package repositories

import (
	"context"
	"strings"

	"github.com/okiww/billing-loan-system/internal/billing_config/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
)

type billingConfigRepository struct {
	*mysql.DBMySQL
}

// GetBillingConfigByName get a BillingConfig by name
func (repo *billingConfigRepository) GetBillingConfigByName(ctx context.Context, name string) (*models.BillingConfig, error) {

	billingConfig := &models.BillingConfig{}
	query := "SELECT * FROM billing_configs WHERE name = ?"
	err := repo.DB.GetContext(ctx, billingConfig, query, strings.ToLower(name))
	if err != nil {
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
