package models

import "time"

// LoanBillModel represents the `loan_bills` table
type LoanBillModel struct {
	ID                 int       `db:"id"`
	LoanID             int64     `db:"loan_id"`
	BillingDate        time.Time `db:"billing_date"`
	BillingAmount      int32     `db:"billing_amount"`       // is the original bill amount
	BillingTotalAmount int32     `db:"billing_total_amount"` // is total that users should pay
	BillingNumber      int       `db:"billing_number"`
	Status             string    `db:"status"` // e.g. 'PENDING', 'PAID', 'OVERDUE'
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
