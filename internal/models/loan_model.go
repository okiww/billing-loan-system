package models

import "time"

type LoanModel struct {
	ID                 int64     `db:"id"`
	UserID             int64     `db:"user_id"`
	Name               string    `db:"name"`
	LoanAmount         int64     `db:"loan_amount"`
	InterestPercentage float64   `db:"interest_percentage"`
	Status             string    `db:"status"`
	StartDate          time.Time `db:"start_date"`
	DueDate            time.Time `db:"due_date"`
	LoanTermsPerWeek   int       `db:"loan_terms_per_week"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
