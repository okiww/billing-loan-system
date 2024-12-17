package models

import "time"

type Payment struct {
	ID         int        `db:"id"`
	UserID     int        `db:"user_id"`
	LoanID     int        `db:"loan_id"`
	LoanBillID int        `db:"loan_bill_id"`
	Amount     int        `db:"amount"`
	Status     string     `db:"status"`
	Note       *string    `db:"note"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

const (
	StatusPending   = "PENDING"
	StatusProcess   = "PROCESS"
	StatusCompleted = "COMPLETED"
	StatusFailed    = "FAILED"

	Note_Complete                 = "Payment Completed"
	Note_Failed_With_ERROR_SYSTEM = "Failed process payment, please try again"
)
