package dto

import (
	"time"

	"github.com/okiww/billing-loan-system/pkg/errors"
)

type LoanRequest struct {
	UserID     int    `json:"user_id"`
	Name       string `json:"name"`
	LoanAmount int32  `json:"loan_amount"`
}

type LoanResponse struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	Name               string    `json:"name"`
	LoanAmount         int       `json:"loan_amount"`
	InterestPercentage float64   `json:"interest_percentage"`
	Status             string    `json:"status"` // ACTIVE, DELINQUENT, CLOSED, PENDING
	StartDate          time.Time `json:"start_date"`
	DueDate            time.Time `json:"due_date"`
	LoanTermsPerWeek   int       `json:"loan_terms_per_week"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (r *LoanRequest) Validate() error {
	// Check UserID
	if r.UserID <= 0 {
		return errors.New("user_id must be greater than 0")
	}

	// Check Name
	if len(r.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	// Check LoanAmount
	if r.LoanAmount <= 0 {
		return errors.New("loan_amount must be greater than 0")
	}

	return nil
}
