package repositories

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/okiww/billing-loan-system/internal/loan/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestGetLoanByID(t *testing.T) {
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
func TestCreateLoan(t *testing.T) {
	// Table-driven test cases
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanRepository(mockDB)

	mockStartDate := time.Date(2024, 12, 16, 10, 0, 0, 0, time.UTC)
	mockDueDate := mockStartDate.Add(30 * 24 * time.Hour)
	type args struct {
		loan *models.LoanModel
	}
	tests := []struct {
		name    string
		s       LoanRepositoryInterface
		args    args
		want    int64
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Loan Created",
			s:    repo,
			args: args{
				loan: &models.LoanModel{
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
			},
			want:    1, // Expected ID of the created loan
			wantErr: false,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO loans`)).
					WithArgs(
						a.loan.UserID,
						a.loan.Name,
						a.loan.LoanAmount,
						a.loan.LoanTotalAmount,
						a.loan.OutstandingAmount,
						a.loan.InterestPercentage,
						a.loan.Status,
						a.loan.StartDate,
						a.loan.DueDate,
						a.loan.LoanTermsPerWeek,
					).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate success, returning ID 1
			},
		},
		{
			name: "Database Error",
			s:    repo,
			args: args{
				loan: &models.LoanModel{
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
			},
			want:    0, // Expected ID is 0 due to error
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO loans`)).
					WithArgs(
						a.loan.UserID,
						a.loan.Name,
						a.loan.LoanAmount,
						a.loan.LoanTotalAmount,
						a.loan.OutstandingAmount,
						a.loan.InterestPercentage,
						a.loan.Status,
						a.loan.StartDate,
						a.loan.DueDate,
						a.loan.LoanTermsPerWeek,
					).
					WillReturnError(errors.New("db error")) // Simulate a DB error
			},
		},
		{
			name: "Validation Error - Missing Required Field",
			s:    repo,
			args: args{
				loan: &models.LoanModel{
					UserID:             123,
					Name:               "",
					LoanAmount:         1000,
					LoanTotalAmount:    1100,
					OutstandingAmount:  1100,
					InterestPercentage: 10,
					Status:             "ACTIVE",
					StartDate:          mockStartDate,
					DueDate:            mockDueDate,
					LoanTermsPerWeek:   4,
				},
			},
			want:    0, // Expected ID is 0 due to validation error
			wantErr: true,
			mock:    func(a args) {}, // No database interaction due to validation failure
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := tt.s.CreateLoan(context.Background(), tt.args.loan)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLoan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("CreateLoan() got = %v, want %v", got, tt.want)
			}

			// Validate expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}
func TestFetchActiveLoan(t *testing.T) {
	// Table-driven test cases
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanRepository(mockDB)

	// Test cases
	tests := []struct {
		name    string
		s       LoanRepositoryInterface
		want    []models.LoanModel
		wantErr bool
		mock    func()
	}{
		{
			name: "Success - Active Loans Found",
			s:    repo,
			want: []models.LoanModel{
				{
					ID:                 1,
					UserID:             123,
					Name:               "Test Loan 1",
					LoanAmount:         1000,
					LoanTotalAmount:    1100,
					OutstandingAmount:  1100,
					InterestPercentage: 10,
					Status:             "ACTIVE",
				},
				{
					ID:                 2,
					UserID:             124,
					Name:               "Test Loan 2",
					LoanAmount:         2000,
					LoanTotalAmount:    2200,
					OutstandingAmount:  2200,
					InterestPercentage: 12,
					Status:             "ACTIVE",
				},
			},
			wantErr: false,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, user_id, name, loan_amount, loan_total_amount, outstanding_amount, 
		       			interest_percentage, status, start_date, due_date, loan_terms_per_week
					FROM loans
					WHERE status = 'ACTIVE'
				`)).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "name", "loan_amount", "loan_total_amount", "outstanding_amount",
						"interest_percentage", "status",
					}).
						AddRow(1, 123, "Test Loan 1", 1000, 1100, 1100, 10, "ACTIVE").
						AddRow(2, 124, "Test Loan 2", 2000, 2200, 2200, 12, "ACTIVE"),
					)
			},
		},
		{
			name:    "No Active Loans Found",
			s:       repo,
			want:    nil,
			wantErr: false,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, user_id, name, loan_amount, loan_total_amount, outstanding_amount, 
		       			interest_percentage, status, start_date, due_date, loan_terms_per_week
					FROM loans
					WHERE status = 'ACTIVE'
				`)).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "name", "loan_amount", "loan_total_amount", "outstanding_amount",
						"interest_percentage", "status",
					}))
			},
		},
		{
			name:    "Database Error",
			s:       repo,
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, user_id, name, loan_amount, loan_total_amount, outstanding_amount, 
		       			interest_percentage, status, start_date, due_date, loan_terms_per_week
					FROM loans
					WHERE status = 'ACTIVE'
				`)).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := tt.s.FetchActiveLoan(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchActiveLoan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchActiveLoan() got = %v, want %v", got, tt.want)
			}

			// Validate expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}
