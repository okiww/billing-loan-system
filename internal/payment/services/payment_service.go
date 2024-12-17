package services

import (
	"context"
	"github.com/okiww/billing-loan-system/internal/dto"
	loanRepo "github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/internal/payment/repositories"
	"github.com/okiww/billing-loan-system/pkg/errors"
)

type paymentService struct {
	paymentRepo repositories.PaymentRepositoryInterface
	loanRepo    loanRepo.LoanRepositoryInterface
}

// MakePayment is for initial payment
func (p *paymentService) MakePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) error {
	// Optionally, you can validate or process business logic here (e.g., check if payment is valid)
	if paymentRequest.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	// Insert the payment into the database
	return p.paymentRepo.Create(ctx, &models.Payment{
		UserID:     paymentRequest.UserID,
		LoanID:     paymentRequest.LoanID,
		LoanBillID: paymentRequest.LoanBillID,
		Amount:     paymentRequest.Amount,
		Status:     models.StatusPending,
	})
}

// ProcessUpdatePayment is for update payment via subscriber
func (p *paymentService) ProcessUpdatePayment(ctx context.Context, payment models.Payment) error {
	//TODO implement me
	// 1. Update Payment to PROCESS
	err := p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusProcess, "")
	if err != nil {
		return err
	}

	// IN TX
	// 2. Update Loan Bills to PAID
	// 3. Check if it is last bill of the loan, update Loan status to CLOSED
	err = p.loanRepo.UpdateLoanAndLoanBillsInTx(ctx, payment.LoanID, payment.LoanBillID, payment.Amount)
	if err != nil {
		// if error, update payment to failed
		err := p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusFailed, models.Note_Failed_With_ERROR_SYSTEM)
		if err != nil {
			return err
		}
		return err
	}
	// DONE TX
	// 4. Update Payment to Success
	err = p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusCompleted, "")
	if err != nil {
		return err
	}
	return nil
}

type PaymentServiceInterface interface {
	MakePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) error
	ProcessUpdatePayment(ctx context.Context, request models.Payment) error
}

func NewPaymentService(paymentRepo repositories.PaymentRepositoryInterface, loanRepo loanRepo.LoanRepositoryInterface) PaymentServiceInterface {
	return &paymentService{
		paymentRepo: paymentRepo,
		loanRepo:    loanRepo,
	}
}
