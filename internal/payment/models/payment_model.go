package models

type Payment struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	LoanID     int    `json:"loan_id"`
	LoanBillID int    `json:"loan_bill_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
}

const (
	StatusPending   = "PENDING"
	StatusProcess   = "PROCESS"
	StatusCompleted = "COMPLETED"
	StatusFailed    = "FAILED"

	Note_Complete                 = "Payment Completed"
	Note_Failed_With_ERROR_SYSTEM = "Failed process payment, please try again"
)
