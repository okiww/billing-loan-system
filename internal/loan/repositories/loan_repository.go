package repositories

import (
	"github.com/okiww/billing-loan-system/internal/models"
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
	return loan, nil
}

// CreateLoan inserts a new loan into the database
func (l *loanRepository) CreateLoan(loan *models.LoanModel) error {
	query := `INSERT INTO loans (user_id, name, loan_amount, interest_percentage, status, start_date, due_date, loan_terms_per_week)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := l.DB.Exec(query, loan.UserID, loan.Name, loan.LoanAmount, loan.InterestPercentage, loan.Status, loan.StartDate, loan.DueDate, loan.LoanTermsPerWeek)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"dataModel": loan,
		}).Error("error when save loan to db")
		return err
	}
	return err
}

type LoanRepositoryInterface interface {
	GetLoanByID(id int64) (*models.LoanModel, error)
	CreateLoan(loan *models.LoanModel) error
}

func NewLoanRepository(db *mysql.DBMySQL) LoanRepositoryInterface {
	return &loanRepository{
		db,
	}
}
