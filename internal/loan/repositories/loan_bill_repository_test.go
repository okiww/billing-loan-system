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

func TestUpdateLoanBillStatuses(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create the repository with the mocked DB
	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanBillRepository(mockDB)

	// Test cases
	tests := []struct {
		name    string
		s       LoanBillRepositoryInterface
		mock    func() // Setup the mock expectations
		wantErr bool   // Expected error state
	}{
		{
			name: "Success - Statuses Updated",
			s:    repo,
			mock: func() {
				// Mock the database query and its result
				mock.ExpectExec(`UPDATE loan_bills`).
					WithArgs().
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate a successful update
			},
			wantErr: false,
		},
		{
			name: "Failure - Database Error",
			s:    repo,
			mock: func() {
				// Mock the database query and simulate an error
				mock.ExpectExec(`UPDATE loan_bills`).
					WithArgs().
					WillReturnError(errors.New("db error")) // Simulate a database error
			},
			wantErr: true,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mock()

			// Call the method
			err := tt.s.UpdateLoanBillStatuses(context.Background())

			// Check if the error state matches the expected result
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLoanBillStatuses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate expectations
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("mock expectations were not met: %v", err)
			}
		})
	}
}

func TestGetTotalLoanBillOverdueByLoanID(t *testing.T) {
	// Initialize mock database and repository
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanBillRepository(mockDB)

	type args struct {
		id int32
	}
	tests := []struct {
		name    string
		repo    LoanBillRepositoryInterface
		args    args
		want    int
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Overdue Loan Bills Found",
			repo: repo,
			args: args{
				id: 1,
			},
			want:    3,
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT COUNT(lb.id) AS overdue_count
					 FROM loan_bills lb
					 JOIN loans l ON lb.loan_id = l.id
					 WHERE l.id = ? AND l.status = 'ACTIVE' AND lb.status = 'OVERDUE'`)).
					WithArgs(a.id).
					WillReturnRows(sqlmock.NewRows([]string{"overdue_count"}).AddRow(3))
			},
		},
		{
			name: "Success - No Overdue Loan Bills",
			repo: repo,
			args: args{
				id: 2,
			},
			want:    0,
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT COUNT(lb.id) AS overdue_count
					 FROM loan_bills lb
					 JOIN loans l ON lb.loan_id = l.id
					 WHERE l.id = ? AND l.status = 'ACTIVE' AND lb.status = 'OVERDUE'`)).
					WithArgs(a.id).
					WillReturnRows(sqlmock.NewRows([]string{"overdue_count"}).AddRow(0))
			},
		},
		{
			name: "Database Error",
			repo: repo,
			args: args{
				id: 3,
			},
			want:    0,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT COUNT(lb.id) AS overdue_count
					 FROM loan_bills lb
					 JOIN loans l ON lb.loan_id = l.id
					 WHERE l.id = ? AND l.status = 'ACTIVE' AND lb.status = 'OVERDUE'`)).
					WithArgs(a.id).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := tt.repo.GetTotalLoanBillOverdueByLoanID(context.Background(), tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalLoanBillOverdueByLoanID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetTotalLoanBillOverdueByLoanID() got = %v, want %v", got, tt.want)
			}

			// Validate mock expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestGetLoanBillsByLoanID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanBillRepository(mockDB)

	type args struct {
		loanID int
	}

	mockStartDate := time.Date(2024, 12, 16, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		repo    LoanBillRepositoryInterface
		args    args
		want    []models.LoanBillModel
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success",
			repo: repo,
			args: args{
				loanID: 1,
			},
			want: []models.LoanBillModel{
				{
					ID:                 1,
					LoanID:             1,
					BillingDate:        mockStartDate,
					BillingAmount:      1000,
					BillingTotalAmount: 1200,
					BillingNumber:      1,
					Status:             "BILLED",
					CreatedAt:          mockStartDate,
					UpdatedAt:          mockStartDate,
				},
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT id, loan_id, billing_date, billing_amount, billing_total_amount, 
							   billing_number, status, created_at, updated_at 
						FROM loan_bills
						WHERE loan_id = ?
						ORDER by billing_number ASC;
					`)).WithArgs(a.loanID).WillReturnRows(sqlmock.NewRows([]string{
					"id", "loan_id", "billing_date", "billing_amount", "billing_total_amount",
					"billing_number", "status", "created_at", "updated_at",
				}).
					AddRow(1, 1, mockStartDate, 1000, 1200, 1, "BILLED", mockStartDate, mockStartDate),
				)
			},
		},
		{
			name: "Loan Bills Not Found",
			repo: repo,
			args: args{
				loanID: 1,
			},
			want:    nil,
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT id, loan_id, billing_date, billing_amount, billing_total_amount, 
							   billing_number, status, created_at, updated_at 
						FROM loan_bills
						WHERE loan_id = ?
						ORDER by billing_number ASC;
					`)).
					WithArgs(a.loanID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "loan_id", "billing_date", "billing_amount", "billing_total_amount",
						"billing_number", "status", "created_at", "updated_at",
					}))
			},
		},
		{
			name: "Database Error",
			repo: repo,
			args: args{
				loanID: 1,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT id, loan_id, billing_date, billing_amount, billing_total_amount, 
							   billing_number, status, created_at, updated_at 
						FROM loan_bills
						WHERE loan_id = ?
						ORDER by billing_number ASC;
					`)).WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := tt.repo.GetLoanBillsByLoanID(context.Background(), tt.args.loanID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoanBillsByLoanID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLoanBillsByLoanID() got = %v, want %v", got, tt.want)
			}

			// Validate mock expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestGetLoanBillByID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewLoanBillRepository(mockDB)

	type args struct {
		id int
	}

	mockStartDate := time.Date(2024, 12, 16, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		repo    LoanBillRepositoryInterface
		args    args
		want    *models.LoanBillModel
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success",
			repo: repo,
			args: args{
				id: 1,
			},
			want: &models.LoanBillModel{
				ID:                 1,
				LoanID:             1,
				BillingDate:        mockStartDate,
				BillingAmount:      1000,
				BillingTotalAmount: 1200,
				BillingNumber:      1,
				Status:             "BILLED",
				CreatedAt:          mockStartDate,
				UpdatedAt:          mockStartDate,
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT * FROM loan_bills WHERE id = ?
					`)).WithArgs(a.id).WillReturnRows(sqlmock.NewRows([]string{
					"id", "loan_id", "billing_date", "billing_amount", "billing_total_amount",
					"billing_number", "status", "created_at", "updated_at",
				}).
					AddRow(1, 1, mockStartDate, 1000, 1200, 1, "BILLED", mockStartDate, mockStartDate),
				)
			},
		},
		{
			name: "Loan Bill Not Found",
			repo: repo,
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT * FROM loan_bills WHERE id = ?
					`)).
					WithArgs(a.id).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Database Error",
			repo: repo,
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT * FROM loan_bills WHERE id = ?
					`)).WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := tt.repo.GetLoanBillByID(context.Background(), tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoanBillByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLoanBillByID() got = %v, want %v", got, tt.want)
			}

			// Validate mock expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}
