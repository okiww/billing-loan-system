package models

import "time"

type LoanModel struct {
	ID                 int64     `db:"id" json:"id"`
	UserID             int64     `db:"user_id" json:"user_id"`
	Name               string    `db:"name" json:"name"`
	LoanAmount         int32     `db:"loan_amount" json:"loan_amount"`                 // Original loan amount
	LoanTotalAmount    int32     `db:"loan_total_amount" json:"loan_total_amount"`     // Total loan amount with interest
	OutstandingAmount  int32     `db:"outstanding_amount" json:"outstanding_amount"`   // Outstanding amount
	InterestPercentage float64   `db:"interest_percentage" json:"interest_percentage"` // Interest percentage
	Status             string    `db:"status" json:"status"`
	StartDate          time.Time `db:"start_date" json:"start_date"`
	DueDate            time.Time `db:"due_date" json:"due_date"`
	LoanTermsPerWeek   int32     `db:"loan_terms_per_week" json:"loan_terms_per_week"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

type LoanWithBills struct {
	ID                 int64           `db:"id" json:"id"`
	UserID             int64           `db:"user_id" json:"user_id"`
	Name               string          `db:"name" json:"name"`
	LoanAmount         int32           `db:"loan_amount" json:"loan_amount"`                 // Original loan amount
	LoanTotalAmount    int32           `db:"loan_total_amount" json:"loan_total_amount"`     // Total loan amount with interest
	OutstandingAmount  int32           `db:"outstanding_amount" json:"outstanding_amount"`   // Outstanding amount
	InterestPercentage float64         `db:"interest_percentage" json:"interest_percentage"` // Interest percentage
	Status             string          `db:"status" json:"status"`
	StartDate          time.Time       `db:"start_date" json:"start_date"`
	DueDate            time.Time       `db:"due_date" json:"due_date"`
	LoanTermsPerWeek   int32           `db:"loan_terms_per_week" json:"loan_terms_per_week"`
	CreatedAt          time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at" json:"updated_at"`
	LoanBills          []LoanBillModel `json:"loan_bills"`
}

type LoanWithBillsResponse struct {
	Loans []LoanWithBills `json:"loans"`
}

const (
	StatusPending = "PENDING"
	StatusBilled  = "BILLED"
	StatusPaid    = "PAID"

	ConfigInterestPercentage  = "loan_interest_percentage"
	ConfigTermsPerWeek        = "loan_term_per_week"
	DefaultInterestPercentage = 10
	DefaultLoanTermsPerWeek   = 50
)
