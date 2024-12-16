package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/okiww/billing-loan-system/internal/user/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUpdateUserToDelinquent(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer func() {
		err := mock.ExpectationsWereMet()
		assert.NoError(t, err, "Unmet expectations: %v", err)
		db.Close()
	}()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewUserRepository(mockDB)
	ctx := context.Background()

	// Table-driven test cases
	type args struct {
		userID int32
	}
	tests := []struct {
		name    string
		s       UserRepositoryInterface
		args    args
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - User Updated",
			s:    repo,
			args: args{
				userID: 1,
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE user SET is_delinquent = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate 1 row updated
			},
		},
		{
			name: "Database Error",
			s:    repo,
			args: args{
				userID: 99,
			},
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE user SET is_delinquent = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnError(fmt.Errorf("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock setup
			tt.mock(tt.args)

			// Execute the method
			err := tt.s.UpdateUserToDelinquent(ctx, tt.args.userID)

			// Validate the result
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserToDelinquent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer func() {
		// Ensure expectations are met after all tests
		err := mock.ExpectationsWereMet()
		assert.NoError(t, err, "Unmet expectations: %v", err)
		db.Close()
	}()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewUserRepository(mockDB)

	// Table-driven test cases
	type args struct {
		userID int32
	}
	tests := []struct {
		name    string
		s       UserRepositoryInterface
		args    args
		want    *models.UserModel
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - User Found",
			s:    repo,
			args: args{
				userID: 1,
			},
			want: &models.UserModel{
				ID:           1,
				Name:         "John Doe",
				IsDelinquent: false,
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, name, is_delinquent FROM user WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_delinquent"}).
						AddRow(1, "John Doe", false))
			},
		},
		{
			name: "User Not Found",
			s:    repo,
			args: args{
				userID: 99,
			},
			want:    nil,
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, name, is_delinquent FROM user WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Database Error",
			s:    repo,
			args: args{
				userID: 2,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, name, is_delinquent FROM user WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnError(fmt.Errorf("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock setup
			tt.mock(tt.args)

			// Execute the method
			got, err := tt.s.GetUserByID(context.Background(), tt.args.userID)

			// Validate the result
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
