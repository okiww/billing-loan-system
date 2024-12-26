package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/okiww/billing-loan-system/helpers"
	"github.com/okiww/billing-loan-system/internal/loan/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

type loanRepository struct {
	*mysql.DBMySQL
}

// GetLoanStatusByID retrieves a loan by its ID
func (l *loanRepository) GetLoanStatusByID(ctx context.Context, id int64) (*models.LoanModel, error) {
	query := "SELECT status FROM loans WHERE id = ?"
	rows, err := l.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loan *models.LoanModel
	for rows.Next() {
		loan = &models.LoanModel{}
		if err := rows.Scan(&loan.Status); err != nil {
			return nil, err
		}
	}

	// Check for any errors during the row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if loan == nil {
		return nil, fmt.Errorf("no loan found with id %d", id)
	}

	return loan, nil
}

// CreateLoan inserts a new loan into the database
func (l *loanRepository) CreateLoan(ctx context.Context, loan *models.LoanModel) (int64, error) {
	query := `INSERT INTO loans (user_id, name, loan_amount, loan_total_amount, outstanding_amount, interest_percentage, status, start_date, due_date, loan_terms_per_week)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := l.DB.ExecContext(ctx, query, loan.UserID, loan.Name, loan.LoanAmount, loan.LoanTotalAmount, loan.OutstandingAmount, loan.InterestPercentage, loan.Status, loan.StartDate, loan.DueDate, loan.LoanTermsPerWeek)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"dataModel": loan,
		}).Error("error when save to loans table")
		return 0, err
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

// FetchActiveLoan retrieves loans with an ACTIVE status
func (l *loanRepository) FetchActiveLoan(ctx context.Context) ([]models.LoanModel, error) {
	query := `
		SELECT id, user_id, name, loan_amount, loan_total_amount, outstanding_amount, 
		       interest_percentage, status, start_date, due_date, loan_terms_per_week
		FROM loans
		WHERE status = 'ACTIVE'
	`
	var activeLoans []models.LoanModel
	err := l.DB.SelectContext(ctx, &activeLoans, query)
	if err != nil {
		return nil, err
	}
	return activeLoans, nil
}

// UpdateLoanAndLoanBillsInTx update loan bills and loans
func (l *loanRepository) UpdateLoanAndLoanBillsInTx(ctx context.Context, loanID, loanBillID, amount int) error {
	err := l.ExecTx(ctx, l.DB, func(tx *sqlx.Tx) error {
		// Update loan bill to Paid
		err := l.UpdateBilledLoanBillToPaid(ctx,
			tx,
			loanBillID,
		)
		if err != nil {
			logger.GetLogger().Errorf("[LoanRepository][UpdateLoanAndLoanBillsInTx] Error UpdateBilledLoanBillToPaid with err: %v", err)
			return err
		}

		// Update Loan Outstanding and Status
		err = l.UpdateOutStandingAmountAndStatus(ctx,
			tx,
			loanID,
			amount,
		)
		if err != nil {
			logger.GetLogger().Errorf("[LoanRepository][UpdateOutStandingAmountAndStatus] Error UpdateBilledLoanBillToPaid with err: %v", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *loanRepository) UpdateBilledLoanBillToPaid(ctx context.Context, tx *sqlx.Tx, id int) error {
	query := `
		UPDATE loan_bills SET status = ?, updated_at = ? WHERE id = ?
	`
	_, err := tx.ExecContext(ctx, query, models.StatusPaid, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (l *loanRepository) UpdateOutStandingAmountAndStatus(ctx context.Context, tx *sqlx.Tx, id, amount int) error {
	query := `
		UPDATE loans
		SET 
			outstanding_amount = outstanding_amount - ?,
			status = CASE
				WHEN outstanding_amount = 0 THEN ?
				ELSE status
    		END
		WHERE id = ?;
	`
	_, err := tx.ExecContext(ctx, query, amount, models.StatusClosed, id)
	if err != nil {
		return err
	}
	return nil
}

func (l *loanRepository) GetLoanByUserID(ctx context.Context, userID int) ([]models.LoanModel, error) {
	query := `
		SELECT id, user_id, name, loan_amount, loan_total_amount, outstanding_amount, 
		       interest_percentage, status, start_date, due_date, loan_terms_per_week
		FROM loans
		WHERE user_id = ?
	`
	var loans []models.LoanModel
	err := l.DB.SelectContext(ctx, &loans, query, userID)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

type LoanRepositoryInterface interface {
	GetLoanStatusByID(ctx context.Context, id int64) (*models.LoanModel, error)
	CreateLoan(ctx context.Context, loan *models.LoanModel) (int64, error)
	FetchActiveLoan(ctx context.Context) ([]models.LoanModel, error)
	UpdateLoanAndLoanBillsInTx(ctx context.Context, loanID, loanBillID, amount int) error
	UpdateBilledLoanBillToPaid(ctx context.Context, tx *sqlx.Tx, id int) error
	UpdateOutStandingAmountAndStatus(ctx context.Context, tx *sqlx.Tx, id, amount int) error
	GetLoanByUserID(ctx context.Context, userID int) ([]models.LoanModel, error)
}

func NewLoanRepository(db *mysql.DBMySQL) LoanRepositoryInterface {
	if helpers.IsTestEnv() { // Skip singleton in tests
		return &loanRepository{
			db,
		}
	}

	repoLoanLock.Do(func() {
		repoLoan = &loanRepository{
			db,
		}
	})
	return repoLoan
}
