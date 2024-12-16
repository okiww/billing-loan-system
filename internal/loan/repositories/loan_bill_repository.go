package repositories

import (
	"context"

	"github.com/okiww/billing-loan-system/internal/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

type loanBillRepository struct {
	*mysql.DBMySQL
}

// GetLoanByID retrieves a loan by its ID
func (l *loanBillRepository) GetLoanByID(id int64) (*models.LoanModel, error) {
	loan := &models.LoanModel{}
	query := "SELECT * FROM loans WHERE id = ?"
	err := l.DB.Get(loan, query, id)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

// CreateLoanBill inserts a new loan_bills into the database
func (l *loanBillRepository) CreateLoanBill(ctx context.Context, loanBill *models.LoanBillModel) error {
	query := `INSERT INTO loan_bills (loan_id, billing_date, billing_amount, billing_total_amount, billing_number, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := l.DB.ExecContext(ctx, query, loanBill.LoanID, loanBill.BillingDate, loanBill.BillingAmount, loanBill.BillingTotalAmount, loanBill.BillingNumber, loanBill.Status, loanBill.CreatedAt, loanBill.UpdatedAt)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"dataModel": loanBill,
		}).Error("error when save to loan_bills table")
		return err
	}
	return err
}

type LoanBillRepositoryInterface interface {
	CreateLoanBill(ctx context.Context, loanBill *models.LoanBillModel) error
}

func NewLoanBillRepository(db *mysql.DBMySQL) LoanBillRepositoryInterface {
	return &loanBillRepository{
		db,
	}
}
