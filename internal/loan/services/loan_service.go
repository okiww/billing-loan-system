package services

import (
	"context"
	"github.com/okiww/billing-loan-system/internal/dto"
	"github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/models"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
	"time"
)

type loanService struct {
	loanRepo repositories.LoanRepositoryInterface
}

func (l *loanService) CreateLoan(ctx context.Context, request dto.LoanRequest) error {
	logger.Info("[LoanService][CreateLoan]")
	// Create a new loan
	newLoan := &models.LoanModel{
		UserID:             int64(request.UserID),
		Name:               request.Name,
		LoanAmount:         int64(request.LoanAmount),
		InterestPercentage: request.InterestPercentage,
		Status:             request.Status,
		StartDate:          time.Now(),
		DueDate:            time.Now().AddDate(0, 6, 0), // 6 months later
		LoanTermsPerWeek:   50,
	}
	err := l.loanRepo.CreateLoan(newLoan)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"request": request,
		}).Error("error when create user to db")
		return err
	}

	return nil
}

type LoanServiceInterface interface {
	CreateLoan(ctx context.Context, request dto.LoanRequest) error
}

func NewLoanService(loanRepo repositories.LoanRepositoryInterface) LoanServiceInterface {
	return &loanService{
		loanRepo,
	}
}
