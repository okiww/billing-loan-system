package services

import (
	"context"
	"fmt"
	"github.com/okiww/billing-loan-system/internal/loan/models"
	"sync"
	"time"

	"github.com/okiww/billing-loan-system/helpers"
	"github.com/okiww/billing-loan-system/internal/dto"
	"github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

type loanService struct {
	loanRepo     repositories.LoanRepositoryInterface
	loanBillRepo repositories.LoanBillRepositoryInterface
}

func (l *loanService) CreateLoan(ctx context.Context, request dto.LoanRequest) error {
	logger.Info("[LoanService][CreateLoan]")
	// Create a new loan
	loanTotalAmount := int32(float64(request.LoanAmount) + (float64(request.LoanAmount) * 10 / 100))
	newLoan := &models.LoanModel{
		UserID:             int64(request.UserID),
		Name:               request.Name,
		LoanAmount:         request.LoanAmount,
		LoanTotalAmount:    loanTotalAmount,
		OutstandingAmount:  request.LoanAmount,
		InterestPercentage: 10, // TODO Should be get From Config
		Status:             request.Status,
		StartDate:          time.Now(),
		DueDate:            helpers.GenerateLastBillDate(time.Now(), 4),
		LoanTermsPerWeek:   4, // TODO should be get from config
	}

	id, err := l.loanRepo.CreateLoan(ctx, newLoan)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"request": request,
		}).Error("error when create loan to db")
		return err
	}

	// Generate loan bills
	err = l.generateLoanBills(ctx, newLoan, id)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"loan_id": newLoan.ID,
		}).Error("error when generating loan bills")
		return err
	}

	return nil
}

// generateLoanBills generates weekly loan bills based on the loan information
func (l *loanService) generateLoanBills(ctx context.Context, loan *models.LoanModel, id int64) error {
	// Use a wait group to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to handle errors from goroutines
	errChan := make(chan error, loan.LoanTermsPerWeek)

	// Define duration for one week
	startBillDate := loan.StartDate

	// Generate loan bills concurrently
	for week := 1; week <= int(loan.LoanTermsPerWeek); week++ {
		// Increment the wait group counter
		wg.Add(1)
		// calculate weekly amount
		weeklyAmount := loan.LoanAmount / loan.LoanTermsPerWeek
		weeklyTotalAmount := loan.LoanTotalAmount / loan.LoanTermsPerWeek
		// Create the billing date (incremented by one week)
		billingDate := helpers.GetNextMonday(startBillDate)

		go func(week int) {
			defer wg.Done() // Decrement the counter when the goroutine finishes

			// Create a new LoanBill model
			loanBill := &models.LoanBillModel{
				LoanID:             id,
				BillingDate:        billingDate,
				BillingAmount:      weeklyAmount,
				BillingTotalAmount: weeklyTotalAmount,
				BillingNumber:      week,
				Status:             "PENDING", // You can adjust this based on the actual status you want
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}

			// Insert the loan bill into the database
			err := l.loanBillRepo.CreateLoanBill(ctx, loanBill)
			if err != nil {
				errChan <- fmt.Errorf("error creating loan bill for week %d: %v", week, err)
				return
			}
		}(week)
		startBillDate = billingDate
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the error channel
	close(errChan)

	// Check for any errors from the goroutines
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

type LoanServiceInterface interface {
	CreateLoan(ctx context.Context, request dto.LoanRequest) error
}

func NewLoanService(loanRepo repositories.LoanRepositoryInterface, loanBillRepo repositories.LoanBillRepositoryInterface) LoanServiceInterface {
	return &loanService{
		loanRepo,
		loanBillRepo,
	}
}
