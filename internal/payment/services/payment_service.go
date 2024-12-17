package services

import (
	"context"

	"github.com/okiww/billing-loan-system/internal/dto"
	loanRepo "github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/internal/payment/repositories"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/okiww/billing-loan-system/pkg/logger"
)

type paymentService struct {
	paymentRepo repositories.PaymentRepositoryInterface
	loanRepo    loanRepo.LoanRepositoryInterface
}

// MakePayment is for initial payment
func (p *paymentService) MakePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) (*models.Payment, error) {
	// Validation if loan_bills.status = 'BILLED'
	// Validation if loans.status = 'ACTIVE
	// Validation if loan_bills = paymentRequest.Amount

	if paymentRequest.Amount <= 0 {
		return nil, errors.New("payment amount must be greater than zero")
	}

	id, err := p.paymentRepo.Create(ctx, &models.Payment{
		UserID:     paymentRequest.UserID,
		LoanID:     paymentRequest.LoanID,
		LoanBillID: paymentRequest.LoanBillID,
		Amount:     paymentRequest.Amount,
		Status:     models.StatusPending,
	})
	if err != nil {
		logger.GetLogger().Errorf("[PaymentSerivce][MakePayment] Error Create with err: %v", err)
		return nil, err
	}

	payment, err := p.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		logger.GetLogger().Errorf("[PaymentSerivce][MakePayment] Error GetPaymentByID with err: %v", err)
		return nil, err
	}

	// Insert the payment into the database
	return payment, nil
}

// ProcessUpdatePayment is for update payment via subscriber
func (p *paymentService) ProcessUpdatePayment(ctx context.Context, payment models.Payment) error {
	logger.GetLogger().Info("[PaymentSerivce][ProcessUpdatePayment]")
	// 1. Update Payment to PROCESS
	err := p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusProcess, "")
	if err != nil {
		logger.GetLogger().Errorf("[PaymentSerivce][ProcessUpdatePayment] Error UpdatePaymentStatus with err: %v", err)
		return err
	}

	// IN TX
	// 2. Update Loan Bills to PAID
	// 3. Check if it is last bill of the loan, update Loan status to CLOSED
	err = p.loanRepo.UpdateLoanAndLoanBillsInTx(ctx, payment.LoanID, payment.LoanBillID, payment.Amount)
	if err != nil {
		logger.GetLogger().Errorf("[PaymentSerivce][UpdatePaymentStatus] Error UpdateLoanAndLoanBillsInTx with err: %v", err)
		// if error, update payment to failed
		err := p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusFailed, models.Note_Failed_With_ERROR_SYSTEM)
		if err != nil {
			logger.GetLogger().Errorf("[PaymentSerivce][UpdatePaymentStatus] Error UpdatePaymentStatus to Failed with err: %v", err)
			return err
		}
		return err
	}
	// DONE TX
	// 4. Update Payment to Completed
	err = p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusCompleted, "")
	if err != nil {
		logger.GetLogger().Errorf("[PaymentSerivce][UpdatePaymentStatus] Error UpdatePaymentStatus to Failed with err: %v", err)
		return err
	}
	return nil
}

type PaymentServiceInterface interface {
	MakePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) (*models.Payment, error)
	ProcessUpdatePayment(ctx context.Context, request models.Payment) error
}

func NewPaymentService(paymentRepo repositories.PaymentRepositoryInterface, loanRepo loanRepo.LoanRepositoryInterface) PaymentServiceInterface {
	return &paymentService{
		paymentRepo: paymentRepo,
		loanRepo:    loanRepo,
	}
}
