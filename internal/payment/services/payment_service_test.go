package services

import (
	"context"
	"testing"

	loan_mock "github.com/okiww/billing-loan-system/gen/mocks/loan"
	payment_mock "github.com/okiww/billing-loan-system/gen/mocks/payment"

	"github.com/golang/mock/gomock"
	"github.com/okiww/billing-loan-system/internal/dto"
	"github.com/okiww/billing-loan-system/internal/loan/models"
	paymentModel "github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMakePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockPaymentRepo := payment_mock.NewMockPaymentRepositoryInterface(ctrl)
	mockLoanRepo := loan_mock.NewMockLoanRepositoryInterface(ctrl)
	mockLoanBillRepo := loan_mock.NewMockLoanBillRepositoryInterface(ctrl)

	// Create the service instance with mocked repos
	service := NewPaymentService(mockPaymentRepo, mockLoanRepo, mockLoanBillRepo)

	// Test table for MakePayment
	tests := []struct {
		name            string
		paymentRequest  *dto.PaymentRequest
		loanBill        *models.LoanBillModel
		loan            *models.LoanModel
		mockRepoCalls   func()
		expectedErr     error
		wantErr         bool
		expectedPayment *paymentModel.Payment
	}{
		{
			name: "Successful Payment",
			paymentRequest: &dto.PaymentRequest{
				UserID:     1,
				LoanID:     1,
				LoanBillID: 1,
				Amount:     1000,
			},
			loanBill: &models.LoanBillModel{Status: models.StatusBilled, BillingTotalAmount: 1000},
			loan:     &models.LoanModel{Status: models.StatusActive},
			mockRepoCalls: func() {
				mockLoanBillRepo.EXPECT().
					GetLoanBillByID(context.Background(), 1).
					Return(&models.LoanBillModel{Status: models.StatusBilled, BillingTotalAmount: 1000}, nil)
				mockLoanRepo.EXPECT().
					GetLoanByID(context.Background(), int64(1)).
					Return(&models.LoanModel{Status: models.StatusActive}, nil)
				mockPaymentRepo.EXPECT().
					Create(context.Background(), gomock.Any()).
					Return(int32(1), nil)
				mockPaymentRepo.EXPECT().
					GetPaymentByID(context.Background(), int32(1)).
					Return(&paymentModel.Payment{ID: 1, Amount: 1000}, nil)
			},
			expectedErr: nil,
			expectedPayment: &paymentModel.Payment{
				ID:     1,
				Amount: 1000,
			},
			wantErr: false,
		},
		{
			name: "Loan Bill Status Not Billed",
			paymentRequest: &dto.PaymentRequest{
				UserID:     1,
				LoanID:     1,
				LoanBillID: 1,
				Amount:     1000,
			},
			loanBill: &models.LoanBillModel{Status: models.StatusPending, BillingTotalAmount: 1000},
			mockRepoCalls: func() {
				mockLoanBillRepo.EXPECT().
					GetLoanBillByID(context.Background(), 1).
					Return(&models.LoanBillModel{Status: models.StatusPending, BillingTotalAmount: 1000}, nil)
			},
			expectedErr:     errors.New(dto.ErrorLoanBillStatusNotBilled),
			expectedPayment: nil,
			wantErr:         true,
		},
		{
			name: "Loan Not Active",
			paymentRequest: &dto.PaymentRequest{
				UserID:     1,
				LoanID:     1,
				LoanBillID: 1,
				Amount:     1000,
			},
			loanBill: &models.LoanBillModel{Status: models.StatusBilled, BillingTotalAmount: 1000},
			loan:     &models.LoanModel{Status: models.StatusClosed},
			mockRepoCalls: func() {
				mockLoanBillRepo.EXPECT().
					GetLoanBillByID(context.Background(), 1).
					Return(&models.LoanBillModel{Status: models.StatusBilled, BillingTotalAmount: 1000}, nil)
				mockLoanRepo.EXPECT().
					GetLoanByID(context.Background(), int64(1)).
					Return(&models.LoanModel{Status: models.StatusClosed}, nil)
			},
			expectedErr:     errors.New(dto.ErrorLoanIsNotActive),
			expectedPayment: nil,
			wantErr:         true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock repository calls
			tt.mockRepoCalls()

			// Call the service method
			payment, err := service.MakePayment(context.Background(), tt.paymentRequest)

			// Assert the expected results
			if tt.wantErr {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
			assert.Equal(t, tt.expectedPayment, payment)
		})
	}
}

func TestProcessUpdatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockPaymentRepo := payment_mock.NewMockPaymentRepositoryInterface(ctrl)
	mockLoanRepo := loan_mock.NewMockLoanRepositoryInterface(ctrl)
	mockLoanBillRepo := loan_mock.NewMockLoanBillRepositoryInterface(ctrl)

	// Create the service instance with mocked repos
	service := NewPaymentService(mockPaymentRepo, mockLoanRepo, mockLoanBillRepo)

	// Test table for ProcessUpdatePayment
	tests := []struct {
		name          string
		payment       paymentModel.Payment
		mockRepoCalls func()
		expectedErr   error
		wantErr       bool
	}{
		{
			name: "Successful Payment Update",
			payment: paymentModel.Payment{
				ID:         1,
				LoanID:     1,
				LoanBillID: 1,
				Amount:     1000,
				Status:     models.StatusPending,
			},
			mockRepoCalls: func() {
				mockPaymentRepo.EXPECT().
					UpdatePaymentStatus(context.Background(), int32(1), paymentModel.StatusProcess, "").
					Return(nil)
				mockLoanRepo.EXPECT().
					UpdateLoanAndLoanBillsInTx(context.Background(), 1, 1, 1000).
					Return(nil)
				mockPaymentRepo.EXPECT().
					UpdatePaymentStatus(context.Background(), int32(1), paymentModel.StatusCompleted, "").
					Return(nil)
			},
			expectedErr: nil,
			wantErr:     false,
		},
		{
			name: "Failed Payment Update",
			payment: paymentModel.Payment{
				ID:         1,
				LoanID:     1,
				LoanBillID: 1,
				Amount:     1000,
				Status:     models.StatusPending,
			},
			mockRepoCalls: func() {
				mockPaymentRepo.EXPECT().
					UpdatePaymentStatus(context.Background(), int32(1), paymentModel.StatusProcess, "").
					Return(nil)

				// Simulate the error in UpdateLoanAndLoanBillsInTx
				mockLoanRepo.EXPECT().
					UpdateLoanAndLoanBillsInTx(context.Background(), 1, 1, 1000).
					Return(errors.New("some error"))

				// Simulate the second call to UpdatePaymentStatus with StatusFailed
				mockPaymentRepo.EXPECT().
					UpdatePaymentStatus(context.Background(), int32(1), paymentModel.StatusFailed, paymentModel.Note_Failed_With_ERROR_SYSTEM).
					Return(nil)
			},
			expectedErr: errors.New("some error"),
			wantErr:     true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock repository calls
			tt.mockRepoCalls()

			// Call the service method
			err := service.ProcessUpdatePayment(context.Background(), tt.payment)
			// Assert the expected results
			if tt.wantErr {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
