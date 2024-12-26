package repositories

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/okiww/billing-loan-system/internal/payment/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewPaymentRepository(mockDB)

	type args struct {
		payment *models.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantID  int32
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Payment Created",
			args: args{
				payment: &models.Payment{
					UserID:     1,
					LoanID:     2,
					LoanBillID: 3,
					Amount:     5000,
					Status:     "PAID",
				},
			},
			wantID:  1,
			wantErr: false,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO payments (user_id, loan_id, loan_bill_id, amount, status, created_at) VALUES (?, ?, ?, ?, ?, ?)")).
					WithArgs(a.payment.UserID, a.payment.LoanID, a.payment.LoanBillID, a.payment.Amount, a.payment.Status, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Database Error",
			args: args{
				payment: &models.Payment{
					UserID:     1,
					LoanID:     2,
					LoanBillID: 3,
					Amount:     5000,
					Status:     "PAID",
				},
			},
			wantID:  0,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO payments (user_id, loan_id, loan_bill_id, amount, status, created_at) VALUES (?, ?, ?, ?, ?, ?)")).
					WithArgs(a.payment.UserID, a.payment.LoanID, a.payment.LoanBillID, a.payment.Amount, a.payment.Status, sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			id, err := repo.Create(context.Background(), tt.args.payment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id != tt.wantID {
				t.Errorf("Create() got = %v, want %v", id, tt.wantID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePaymentStatus(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewPaymentRepository(mockDB)

	type args struct {
		id     int32
		status string
		note   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Payment Status Updated",
			args: args{
				id:     1,
				status: "COMPLETED",
				note:   "Payment received",
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE payments SET status = ?, updated_at = ?, note = ? WHERE id = ?")).
					WithArgs(a.status, sqlmock.AnyArg(), a.note, a.id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Failure - Database Error",
			args: args{
				id:     1,
				status: "COMPLETED",
				note:   "Payment received",
			},
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE payments SET status = ?, updated_at = ?, note = ? WHERE id = ?")).
					WithArgs(a.status, sqlmock.AnyArg(), a.note, a.id).
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			err := repo.UpdatePaymentStatus(context.Background(), tt.args.id, tt.args.status, tt.args.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePaymentStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPaymentByID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewPaymentRepository(mockDB)

	type args struct {
		id int32
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Payment
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - Payment Found",
			args: args{id: 1},
			want: &models.Payment{
				ID:         1,
				UserID:     1,
				LoanID:     2,
				LoanBillID: 3,
				Amount:     5000,
				Status:     "PAID",
				CreatedAt:  time.Now(),
			},
			wantErr: false,
			mock: func(a args) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "loan_id", "loan_bill_id", "amount", "status", "created_at"}).
					AddRow(1, 1, 2, 3, 5000, "PAID", time.Now())

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, loan_id, loan_bill_id, amount, status, created_at, updated_at, note FROM payments WHERE id = ?")).
					WithArgs(a.id).
					WillReturnRows(rows)
			},
		},
		{
			name:    "Failure - Payment Not Found",
			args:    args{id: 1},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, loan_id, loan_bill_id, amount, status, created_at, updated_at, note FROM payments WHERE id = ?")).
					WithArgs(a.id).
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			got, err := repo.GetPaymentByID(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaymentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.ID, got.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
