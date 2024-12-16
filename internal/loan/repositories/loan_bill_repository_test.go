package repositories

import (
	"context"
	"database/sql"
	"github.com/okiww/billing-loan-system/internal/loan/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestGetLoanBillByID(t *testing.T) {
	// Table-driven test cases
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanRepository(mockDB)

	mockStartDate := time.Date(2024, 12, 16, 10, 0, 0, 0, time.UTC)
	mockDueDate := mockStartDate.Add(30 * 24 * time.Hour)
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		s       LoanRepositoryInterface
		args    args
		want    *models.LoanModel
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Loan Found",
			s:    repo,
			args: args{
				id: 1,
			},
			want: &models.LoanModel{
				ID:                 1,
				UserID:             123,
				Name:               "Test Loan",
				LoanAmount:         1000,
				LoanTotalAmount:    1100,
				OutstandingAmount:  1100,
				InterestPercentage: 10,
				Status:             "ACTIVE",
				StartDate:          mockStartDate,
				DueDate:            mockDueDate,
				LoanTermsPerWeek:   4,
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM loans WHERE id = ?`)).
					WithArgs(a.id).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "name", "loan_amount", "loan_total_amount", "outstanding_amount", "interest_percentage", "status", "start_date", "due_date", "loan_terms_per_week",
					}).AddRow(
						1, 123, "Test Loan", 1000, 1100, 1100, 10, "ACTIVE", mockStartDate, mockDueDate, 4,
					))
			},
		},
		{
			name: "Loan Not Found",
			s:    repo,
			args: args{
				id: 99,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM loans WHERE id = ?`)).
					WithArgs(a.id).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Database Error",
			s:    repo,
			args: args{
				id: 2,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM loans WHERE id = ?`)).
					WithArgs(a.id).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := tt.s.GetLoanByID(tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoanByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLoanByID() got = %v, want %v", got, tt.want)
			}

			// Validate expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestCreateLoanBill(t *testing.T) {
	// Table-driven test cases
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanBillRepository(mockDB)

	mockBillingDate := time.Date(2024, 12, 16, 10, 0, 0, 0, time.UTC)
	mockCreatedAt := mockBillingDate
	mockUpdatedAt := mockBillingDate.Add(24 * time.Hour)
	type args struct {
		loanBill *models.LoanBillModel
	}
	tests := []struct {
		name    string
		s       LoanBillRepositoryInterface
		args    args
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Loan Bill Created",
			s:    repo,
			args: args{
				loanBill: &models.LoanBillModel{
					LoanID:             1,
					BillingDate:        mockBillingDate,
					BillingAmount:      1000,
					BillingTotalAmount: 1100,
					BillingNumber:      1,
					Status:             "PAID",
					CreatedAt:          mockCreatedAt,
					UpdatedAt:          mockUpdatedAt,
				},
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO loan_bills`)).
					WithArgs(
						a.loanBill.LoanID,
						a.loanBill.BillingDate,
						a.loanBill.BillingAmount,
						a.loanBill.BillingTotalAmount,
						a.loanBill.BillingNumber,
						a.loanBill.Status,
						a.loanBill.CreatedAt,
						a.loanBill.UpdatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate success
			},
		},
		{
			name: "Database Error",
			s:    repo,
			args: args{
				loanBill: &models.LoanBillModel{
					LoanID:             1,
					BillingDate:        mockBillingDate,
					BillingAmount:      1000,
					BillingTotalAmount: 1100,
					BillingNumber:      1,
					Status:             "PAID",
					CreatedAt:          mockCreatedAt,
					UpdatedAt:          mockUpdatedAt,
				},
			},
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO loan_bills`)).
					WithArgs(
						a.loanBill.LoanID,
						a.loanBill.BillingDate,
						a.loanBill.BillingAmount,
						a.loanBill.BillingTotalAmount,
						a.loanBill.BillingNumber,
						a.loanBill.Status,
						a.loanBill.CreatedAt,
						a.loanBill.UpdatedAt,
					).
					WillReturnError(errors.New("db error")) // Simulate DB error
			},
		},
		{
			name: "Missing Required Field - Billing Number",
			s:    repo,
			args: args{
				loanBill: &models.LoanBillModel{
					LoanID:             1,
					BillingDate:        mockBillingDate,
					BillingAmount:      1000,
					BillingTotalAmount: 1100,
					BillingNumber:      1, // Missing billing number
					Status:             "PAID",
					CreatedAt:          mockCreatedAt,
					UpdatedAt:          mockUpdatedAt,
				},
			},
			wantErr: true,
			mock:    func(a args) {}, // No database interaction due to validation failure
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			err := tt.s.CreateLoanBill(context.Background(), tt.args.loanBill)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLoanBill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}
