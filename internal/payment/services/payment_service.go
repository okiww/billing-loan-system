package services

import (
	"context"
	"github.com/okiww/billing-loan-system/internal/dto"
	"time"

	"github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/internal/payment/repositories"
	"github.com/okiww/billing-loan-system/pkg/errors"
)

type paymentService struct {
	paymentRepo repositories.PaymentRepositoryInterface
}

func (p *paymentService) CreatePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) error {
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
		CreatedAt:  time.Now(),
	})
}

type PaymentServiceInterface interface {
	CreatePayment(ctx context.Context, paymentRequest *dto.PaymentRequest) error
}

func NewPaymentService(paymentRepo repositories.PaymentRepositoryInterface) PaymentServiceInterface {
	return &paymentService{
		paymentRepo: paymentRepo,
	}
}
