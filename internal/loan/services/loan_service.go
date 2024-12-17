package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/okiww/billing-loan-system/internal/loan/models"

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
	logger.GetLogger().Info("[LoanService][CreateLoan]")
	// Create a new loan
	loanTotalAmount := int32(float64(request.LoanAmount) + (float64(request.LoanAmount) * 10 / 100))
	newLoan := &models.LoanModel{
		UserID:             int64(request.UserID),
		Name:               request.Name,
		LoanAmount:         request.LoanAmount,
		LoanTotalAmount:    loanTotalAmount,
		OutstandingAmount:  loanTotalAmount,
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

func (l *loanService) UpdateLoanBill(ctx context.Context) error {
	logger.GetLogger().Info("[LoanService][UpdateLoanBill]")
	// this is for update loan bill from pending to billed
	err := l.loanBillRepo.UpdateLoanBillStatuses(ctx)
	if err != nil {
		logger.GetLogger().Error("[LoanService][UpdateLoanBill] Error when update loan bill statuses")
		return err
	}
	// this is for update loan that already overdue but still billed to overdue
	// TODO calculate user is_delinquent(if has more 2 overdue loan_bills)
	return nil
}

func (l *loanService) GetAllActiveLoan(ctx context.Context) ([]models.LoanModel, error) {
	logger.GetLogger().Info("[LoanService][GetAllActiveLoan]")
	loans, err := l.loanRepo.FetchActiveLoan(ctx)
	if err != nil {
		logger.GetLogger().Errorf("[LoanService][GetAllActiveLoan] Error when fetch active loans with err: %v", err)
		return []models.LoanModel{}, err
	}

	return loans, nil
}

func (l *loanService) CountLoanBillOverdueStatusesByID(ctx context.Context, id int32) (int32, error) {
	logger.GetLogger().Info("[LoanService][CountLoanBillOverdueStatuses]")
	total, err := l.loanBillRepo.GetTotalLoanBillOverdueByLoanID(ctx, id)
	if err != nil {
		logger.GetLogger().Errorf("[LoanService][CountLoanBillOverdueStatuses] Error when get total loan bill overdue with err: %v", err)
		return 0, err
	}

	return int32(total), nil
}

// generateLoanBills generates weekly loan bills based on the loan information
func (l *loanService) generateLoanBills(ctx context.Context, loan *models.LoanModel, id int64) error {
	logger.GetLogger().Info("[LoanService][generateLoanBills] Start")
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
				Status:             models.StatusPending, // You can adjust this based on the actual status you want
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
	logger.GetLogger().Info("[LoanService][generateLoanBills] Done")
	return nil
}

type LoanServiceInterface interface {
	GetAllActiveLoan(ctx context.Context) ([]models.LoanModel, error)
	CreateLoan(ctx context.Context, request dto.LoanRequest) error
	UpdateLoanBill(ctx context.Context) error
	CountLoanBillOverdueStatusesByID(ctx context.Context, id int32) (int32, error)
}

func NewLoanService(loanRepo repositories.LoanRepositoryInterface, loanBillRepo repositories.LoanBillRepositoryInterface) LoanServiceInterface {
	return &loanService{
		loanRepo,
		loanBillRepo,
	}
}
