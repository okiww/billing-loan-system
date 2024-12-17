package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	billing_config_mock "github.com/okiww/billing-loan-system/gen/mocks/billing_config"
	models2 "github.com/okiww/billing-loan-system/internal/billing_config/models"
	"github.com/okiww/billing-loan-system/internal/dto"

	"github.com/golang/mock/gomock"
	loan_mock "github.com/okiww/billing-loan-system/gen/mocks/loan"

	"github.com/okiww/billing-loan-system/internal/loan/models"
)

func TestCreateLoan(t *testing.T) {
	// Create a new mock controller and defer its cleanup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the loan and loanBill repositories
	mockLoanRepo := loan_mock.NewMockLoanRepositoryInterface(ctrl)
	mockLoanBillRepo := loan_mock.NewMockLoanBillRepositoryInterface(ctrl)
	mockBillingConfig := billing_config_mock.NewMockBillingConfigRepositoryInterface(ctrl)

	// Initialize the loan service
	loanService := NewLoanService(mockLoanRepo, mockLoanBillRepo, mockBillingConfig)

	// Define test cases using a table-driven approach
	tests := []struct {
		name    string
		request dto.LoanRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "Success - Create Loan with Configs",
			request: dto.LoanRequest{
				UserID:     1,
				Name:       "John Doe",
				LoanAmount: 10000,
				Status:     "PENDING",
			},
			setup: func() {
				// Mock getting interest percentage config
				mockBillingConfig.EXPECT().
					GetBillingConfigByName(gomock.Any(), gomock.Eq(models.ConfigInterestPercentage)).
					Return(&models2.BillingConfig{
						ID:    1,
						Name:  models.ConfigInterestPercentage,
						Value: `{"is_active":true,"value":10}`, // Example JSON value
					}, nil).
					Times(1)

				// Mock getting loan terms per week config
				mockBillingConfig.EXPECT().
					GetBillingConfigByName(gomock.Any(), gomock.Eq(models.ConfigTermsPerWeek)).
					Return(&models2.BillingConfig{
						ID:    2,
						Name:  models.ConfigTermsPerWeek,
						Value: `{"is_active":true,"value":5}`, // Example JSON value
					}, nil).
					Times(1)

				// Mock loan repository create loan
				mockLoanRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(int64(1), nil)

				// Mock loan bill repository create loan bills
				mockLoanBillRepo.EXPECT().CreateLoanBill(gomock.Any(), gomock.Any()).Times(5) // 5 bills for 5 weeks
			},
			wantErr: false,
		},
		{
			name: "Error - Loan Creation Failed",
			request: dto.LoanRequest{
				UserID:     1,
				Name:       "John Doe",
				LoanAmount: 10000,
				Status:     "PENDING",
			},
			setup: func() {
				// Mock getting interest percentage config
				mockBillingConfig.EXPECT().
					GetBillingConfigByName(gomock.Any(), gomock.Eq(models.ConfigInterestPercentage)).
					Return(&models2.BillingConfig{
						ID:    1,
						Name:  models.ConfigInterestPercentage,
						Value: `{"is_active":true,"value":10}`, // Example JSON value
					}, nil).
					Times(1)

				// Mock getting loan terms per week config
				mockBillingConfig.EXPECT().
					GetBillingConfigByName(gomock.Any(), gomock.Eq(models.ConfigTermsPerWeek)).
					Return(&models2.BillingConfig{
						ID:    2,
						Name:  models.ConfigTermsPerWeek,
						Value: `{"is_active":true,"value":5}`, // Example JSON value
					}, nil).
					Times(1)

				// Mock loan repository create loan
				mockLoanRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("failed to create loan"))
			},
			wantErr: true,
		},
		{
			name: "Error - Loan Bill Creation Failed",
			request: dto.LoanRequest{
				UserID:     1,
				Name:       "John Doe",
				LoanAmount: 10000,
				Status:     "PENDING",
			},
			setup: func() {
				// Mock getting interest percentage config
				mockBillingConfig.EXPECT().
					GetBillingConfigByName(gomock.Any(), gomock.Eq(models.ConfigInterestPercentage)).
					Return(&models2.BillingConfig{
						ID:    1,
						Name:  models.ConfigInterestPercentage,
						Value: `{"is_active":true,"value":10}`, // Example JSON value
					}, nil).
					Times(1)

				// Mock getting loan terms per week config
				mockBillingConfig.EXPECT().
					GetBillingConfigByName(gomock.Any(), gomock.Eq(models.ConfigTermsPerWeek)).
					Return(&models2.BillingConfig{
						ID:    2,
						Name:  models.ConfigTermsPerWeek,
						Value: `{"is_active":true,"value":5}`, // Example JSON value
					}, nil).
					Times(1)

				// Mock loan repository create loan
				mockLoanRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(int64(1), nil)

				// Mock loan bill repository create loan bills (simulate error for one bill)
				mockLoanBillRepo.EXPECT().CreateLoanBill(gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to create loan bill")).Times(5)
			},
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := loanService.CreateLoan(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLoan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestCreateLoan(t *testing.T) {
//	// Create a new mock controller and defer its clean up
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// Mock the loan and loanBill repositories
//	mockLoanRepo := loan_mock.NewMockLoanRepositoryInterface(ctrl)
//	mockLoanBillRepo := loan_mock.NewMockLoanBillRepositoryInterface(ctrl)
//	mockBillingConfig := billing_config_mock.NewMockBillingConfigRepositoryInterface(ctrl)
//
//	// Initialize the loan service
//	loanService := NewLoanService(mockLoanRepo, mockLoanBillRepo, mockBillingConfig)
//
//	// Define test cases using a table-driven approach
//	tests := []struct {
//		name    string
//		request dto.LoanRequest
//		setup   func()
//		wantErr bool
//	}{
//		{
//			name: "Success - Create Loan",
//			request: dto.LoanRequest{
//				UserID:     1,
//				Name:       "John Doe",
//				LoanAmount: 10000,
//				Status:     "PENDING",
//			},
//			setup: func() {
//				mockLoanRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(int64(1), nil)
//				mockLoanBillRepo.EXPECT().CreateLoanBill(gomock.Any(), gomock.Any()).Times(4) // 4 bills for 4 weeks
//			},
//			wantErr: false,
//		},
//		{
//			name: "Error - Loan Creation Failed",
//			request: dto.LoanRequest{
//				UserID:     1,
//				Name:       "John Doe",
//				LoanAmount: 10000,
//				Status:     "PENDING",
//			},
//			setup: func() {
//				mockLoanRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("failed to create loan"))
//			},
//			wantErr: true,
//		},
//		{
//			name: "Error - Loan Bill Creation Failed",
//			request: dto.LoanRequest{
//				UserID:     1,
//				Name:       "John Doe",
//				LoanAmount: 10000,
//				Status:     "PENDING",
//			},
//			setup: func() {
//				mockLoanRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(int64(1), nil)
//				mockLoanBillRepo.EXPECT().CreateLoanBill(gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to create loan bill")).Times(4).Return(fmt.Errorf("some error"))
//			},
//			wantErr: true,
//		},
//	}
//
//	// Run tests
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setup()
//
//			err := loanService.CreateLoan(context.Background(), tt.request)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("CreateLoan() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestGenerateLoanBills(t *testing.T) {
	// Create a new mock controller and defer its clean up
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the loanBill repository
	mockLoanBillRepo := loan_mock.NewMockLoanBillRepositoryInterface(ctrl)

	// Initialize the loan service
	loanService := &loanService{
		loanBillRepo: mockLoanBillRepo,
	}

	// Define test cases using a table-driven approach
	tests := []struct {
		name    string
		loan    *models.LoanModel
		setup   func()
		wantErr bool
	}{
		{
			name: "Success - Generate Loan Bills",
			loan: &models.LoanModel{
				ID:               1,
				LoanAmount:       10000,
				LoanTotalAmount:  11000,
				LoanTermsPerWeek: 4,
				StartDate:        time.Now(),
			},
			setup: func() {
				// Mock the loanBill repository to return no error for creating bills
				mockLoanBillRepo.EXPECT().CreateLoanBill(gomock.Any(), gomock.Any()).Times(4) // 4 bills for 4 weeks
			},
			wantErr: false,
		},
		{
			name: "Error - Loan Bill Creation Failed",
			loan: &models.LoanModel{
				ID:               1,
				LoanAmount:       10000,
				LoanTotalAmount:  11000,
				LoanTermsPerWeek: 4,
				StartDate:        time.Now(),
			},
			setup: func() {
				// Mock the loanBill repository to return an error for the first bill
				mockLoanBillRepo.EXPECT().CreateLoanBill(gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to create loan bill")).Times(4)
			},
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := loanService.generateLoanBills(context.Background(), tt.loan, tt.loan.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateLoanBills() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
