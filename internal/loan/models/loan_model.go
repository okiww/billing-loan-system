package models

import "time"

type LoanModel struct {
	ID                 int64     `db:"id"`
	UserID             int64     `db:"user_id"`
	Name               string    `db:"name"`
	LoanAmount         int32     `db:"loan_amount"`        // is original amount
	LoanTotalAmount    int32     `db:"loan_total_amount"`  // is total loan amount with interest
	OutstandingAmount  int32     `db:"outstanding_amount"` // is outstanding amount
	InterestPercentage float64   `db:"interest_percentage"`
	Status             string    `db:"status"`
	StartDate          time.Time `db:"start_date"`
	DueDate            time.Time `db:"due_date"`
	LoanTermsPerWeek   int32     `db:"loan_terms_per_week"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

const (
	StatusPending = "PENDING"
	StatusBilled  = "BILLED"
	StatusPaid    = "PAID"
)
