package models

import "time"

// LoanBillModel represents the `loan_bills` table
type LoanBillModel struct {
	ID                 int       `db:"id" json:"id"`
	LoanID             int64     `db:"loan_id" json:"loan_id"`
	BillingDate        time.Time `db:"billing_date" json:"billing_date"`
	BillingAmount      int32     `db:"billing_amount" json:"billing_amount"`             // Original bill amount
	BillingTotalAmount int32     `db:"billing_total_amount" json:"billing_total_amount"` // Total payment amount
	BillingNumber      int       `db:"billing_number" json:"billing_number"`
	Status             string    `db:"status" json:"status"` // e.g., 'PENDING', 'PAID', 'OVERDUE'
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

const (
	StatusActive = "ACTIVE"
	StatusClosed = "CLOSED"
)
