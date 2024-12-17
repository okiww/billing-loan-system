package dto

import "github.com/okiww/billing-loan-system/pkg/errors"

type PaymentRequest struct {
	UserID     int `json:"user_id"`
	LoanID     int `json:"loan_id"`
	Amount     int `json:"amount"`
	LoanBillID int `json:"loan_bill_id"`
}

func (r *PaymentRequest) Validate() error {
	if r.UserID == 0 {
		return errors.New("user_id is required")
	}
	if r.LoanID == 0 {
		return errors.New("loan_id is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

const (
	ErrorLoanBillStatusNotBilled       = "loan bill status is not billed"
	ErrorPaymentAmountNotMatchWithBill = "payment amount not match with the bill"
	ErrorLoanIsNotActive               = "loan is not active"
)
