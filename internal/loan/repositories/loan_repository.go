package repositories

import (
	"context"
	"github.com/okiww/billing-loan-system/internal/loan/models"
	"sync"

	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	repo     LoanRepositoryInterface
	repoLock sync.Once
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

type LoanRepositoryInterface interface {
	GetLoanByID(id int64) (*models.LoanModel, error)
	CreateLoan(ctx context.Context, loan *models.LoanModel) (int64, error)
}

func NewLoanRepository(db *mysql.DBMySQL) LoanRepositoryInterface {
	repoLock.Do(func() {
		repo = &loanRepository{
			db,
		}
	})
	return repo
}
