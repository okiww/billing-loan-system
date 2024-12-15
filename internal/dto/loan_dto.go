package dto

import "time"

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
