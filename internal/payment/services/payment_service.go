package services

import (
	"context"
	loanModel "github.com/okiww/billing-loan-system/internal/loan/models"

	"github.com/okiww/billing-loan-system/internal/dto"
	loanRepo "github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/internal/payment/repositories"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/okiww/billing-loan-system/pkg/logger"
)

type paymentService struct {
	paymentRepo  repositories.PaymentRepositoryInterface
	loanRepo     loanRepo.LoanRepositoryInterface
	loanBillRepo loanRepo.LoanBillRepositoryInterface
}

// MakePayment is for initial payment
func (p *paymentService) MakePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) (*models.Payment, error) {
	logger.GetLogger().Info("[PaymentService][MakePayment]")
	// Validation if loan_bills.status = 'BILLED'
	loanBill, err := p.loanBillRepo.GetLoanBillByID(ctx, int(int64(paymentRequest.LoanBillID)))
	if err != nil {
		return nil, err
	}

	if loanBill.Status != loanModel.StatusBilled {
		return nil, errors.New(dto.ErrorLoanBillStatusNotBilled)
	}

	if loanBill.BillingTotalAmount != int32(paymentRequest.Amount) {
		return nil, errors.New(dto.ErrorPaymentAmountNotMatchWithBill)
	}
	// Validation if loans.status = 'ACTIVE
	loan, err := p.loanRepo.GetLoanByID(ctx, int64(paymentRequest.LoanID))
	if err != nil {
		return nil, err
	}

	if loan.Status != loanModel.StatusActive {
		return nil, errors.New(dto.ErrorLoanIsNotActive)
	}
	// Validation if loan_bills = paymentRequest.Amount

	id, err := p.paymentRepo.Create(ctx, &models.Payment{
		UserID:     paymentRequest.UserID,
		LoanID:     paymentRequest.LoanID,
		LoanBillID: paymentRequest.LoanBillID,
		Amount:     paymentRequest.Amount,
		Status:     models.StatusPending,
	})
	if err != nil {
		logger.GetLogger().Errorf("[PaymentService][MakePayment] Error Create with err: %v", err)
		return nil, err
	}

	payment, err := p.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		logger.GetLogger().Errorf("[PaymentService][MakePayment] Error GetPaymentByID with err: %v", err)
		return nil, err
	}

	// Insert the payment into the database
	return payment, nil
}

// ProcessUpdatePayment is for update payment via subscriber
func (p *paymentService) ProcessUpdatePayment(ctx context.Context, payment models.Payment) error {
	logger.GetLogger().Info("[PaymentService][ProcessUpdatePayment]")
	// 1. Update Payment to PROCESS
	err := p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusProcess, "")
	if err != nil {
		logger.GetLogger().Errorf("[PaymentService][ProcessUpdatePayment] Error UpdatePaymentStatus with err: %v", err)
		return err
	}

	// IN TX
	// 2. Update Loan Bills to PAID
	// 3. Check if it is last bill of the loan, update Loan status to CLOSED
	err = p.loanRepo.UpdateLoanAndLoanBillsInTx(ctx, payment.LoanID, payment.LoanBillID, payment.Amount)
	if err != nil {
		logger.GetLogger().Errorf("[PaymentService][UpdatePaymentStatus] Error UpdateLoanAndLoanBillsInTx with err: %v", err)
		// if error, update payment to failed
		updateErr := p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusFailed, models.Note_Failed_With_ERROR_SYSTEM)
		if updateErr != nil {
			logger.GetLogger().Errorf("[PaymentService][UpdatePaymentStatus] Error UpdatePaymentStatus to Failed with err: %v", err)
			return updateErr
		}
		return err
	}
	// DONE TX
	// 4. Update Payment to Completed
	err = p.paymentRepo.UpdatePaymentStatus(ctx, int32(payment.ID), models.StatusCompleted, "")
	if err != nil {
		logger.GetLogger().Errorf("[PaymentService][UpdatePaymentStatus] Error UpdatePaymentStatus to Failed with err: %v", err)
		return err
	}
	return nil
}

type PaymentServiceInterface interface {
	MakePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) (*models.Payment, error)
	ProcessUpdatePayment(ctx context.Context, request models.Payment) error
}

func NewPaymentService(paymentRepo repositories.PaymentRepositoryInterface, loanRepo loanRepo.LoanRepositoryInterface, loanBillRepo loanRepo.LoanBillRepositoryInterface) PaymentServiceInterface {
	return &paymentService{
		paymentRepo:  paymentRepo,
		loanRepo:     loanRepo,
		loanBillRepo: loanBillRepo,
	}
}
