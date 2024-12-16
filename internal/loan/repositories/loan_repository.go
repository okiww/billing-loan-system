package repositories

import (
	"context"
	"fmt"
	"github.com/okiww/billing-loan-system/helpers"
	"github.com/okiww/billing-loan-system/internal/loan/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

type loanRepository struct {
	*mysql.DBMySQL
}

// GetLoanByID retrieves a loan by its ID
func (l *loanRepository) GetLoanByID(id int64) (*models.LoanModel, error) {
	loan := &models.LoanModel{}
	query := "SELECT * FROM loans WHERE id = ?"
	err := l.DB.Get(loan, query, id)
	if err != nil {
		return nil, err
	}

	fmt.Println(loan)
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

type LoanRepositoryInterface interface {
	GetLoanByID(id int64) (*models.LoanModel, error)
	CreateLoan(ctx context.Context, loan *models.LoanModel) (int64, error)
	FetchActiveLoan(ctx context.Context) ([]models.LoanModel, error)
}

func NewLoanRepository(db *mysql.DBMySQL) LoanRepositoryInterface {
	if helpers.IsTestEnv() { // Skip singleton in tests
		return &loanRepository{
			db,
		}
	}

	repoLock.Do(func() {
		repoLoan = &loanRepository{
			db,
		}
	})
	return repoLoan
}
