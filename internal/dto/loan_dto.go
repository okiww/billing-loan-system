package dto

import (
	"github.com/okiww/billing-loan-system/pkg/errors"
	"time"
)

type LoanRequest struct {
	UserID             int       `json:"user_id"`
	Name               string    `json:"name"`
	LoanAmount         int       `json:"loan_amount"`
	InterestPercentage float64   `json:"interest_percentage"`
	Status             string    `json:"status"` // ACTIVE, DELINQUENT, CLOSED, PENDING
	StartDate          time.Time `json:"start_date"`
	DueDate            time.Time `json:"due_date"`
	LoanTermsPerWeek   int       `json:"loan_terms_per_week"`
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

	// Check InterestPercentage
	if r.InterestPercentage < 0 || r.InterestPercentage > 100 {
		return errors.New("interest_percentage must be between 0 and 100")
	}

	// Check Status
	validStatuses := map[string]bool{
		"ACTIVE":     true,
		"DELINQUENT": true,
		"CLOSED":     true,
		"PENDING":    true,
	}
	if !validStatuses[r.Status] {
		return errors.New("invalid status; must be one of: ACTIVE, DELINQUENT, CLOSED, PENDING")
	}

	// Check StartDate
	if r.StartDate.IsZero() {
		return errors.New("start_date must be specified")
	}

	// Check DueDate
	if r.DueDate.IsZero() {
		return errors.New("due_date must be specified")
	}
	if r.DueDate.Before(r.StartDate) {
		return errors.New("due_date must be after start_date")
	}

	// Check LoanTermsPerWeek
	if r.LoanTermsPerWeek <= 0 {
		return errors.New("loan_terms_per_week must be greater than 0")
	}

	return nil
}
