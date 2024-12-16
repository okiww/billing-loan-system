package repositories

import (
	"context"
	models2 "github.com/okiww/billing-loan-system/internal/loan/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

type loanBillRepository struct {
	*mysql.DBMySQL
}

// CreateLoanBill inserts a new loan_bills into the database
func (l *loanBillRepository) CreateLoanBill(ctx context.Context, loanBill *models2.LoanBillModel) error {
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

func (l *loanBillRepository) UpdateLoanBillStatuses(ctx context.Context) error {
	query := `
		UPDATE loan_bills 
		SET status = CASE
			WHEN billing_date = CURDATE() THEN 'BILLED'
			WHEN billing_date < CURDATE() THEN 'OVERDUE'
		END
		WHERE loan_id IN (
			SELECT id 
			FROM loans 
			WHERE status = 'ACTIVE'
		) 
		AND (billing_date <= CURDATE())
	`

	_, err := l.DB.ExecContext(ctx, query)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to update loan bill statuses")
		return err
	}

	logger.GetLogger().Info("loan bill statuses updated successfully")
	return nil
}

type LoanBillRepositoryInterface interface {
	CreateLoanBill(ctx context.Context, loanBill *models2.LoanBillModel) error
	UpdateLoanBillStatuses(ctx context.Context) error
}

func NewLoanBillRepository(db *mysql.DBMySQL) LoanBillRepositoryInterface {
	return &loanBillRepository{
		db,
	}
}
