package models

import (
	"time"
)

type Payment struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	LoanID     int       `json:"loan_id"`
	LoanBillID int       `json:"loan_bill_id"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

const (
	StatusPending   = "PENDING"
	StatusProcess   = "PROCESS"
	StatusCompleted = "COMPLETED"
)
