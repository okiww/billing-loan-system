package repositories

import (
	"context"
	"database/sql"
	"github.com/okiww/billing-loan-system/internal/user/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"reflect"
	"regexp"
	"testing"
)

func TestUpdateUserToDelinquent(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewUserRepository(mockDB)
	ctx := context.Background()

	// Table-driven test cases
	type args struct {
		userID int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - User Updated",
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
			name: "User Not Found",
			args: args{
				userID: 99,
			},
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE user SET is_delinquent = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnError(errors.New("some error"))
			},
		},
		{
			name: "Database Error",
			args: args{
				userID: 2,
			},
			wantErr: true,
			mock: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE user SET is_delinquent = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			err := repo.UpdateUserToDelinquent(ctx, tt.args.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserToDelinquent() error = %v, wantErr %v", err, tt.wantErr)
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

func TestGetUserByID(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewUserRepository(mockDB)
	ctx := context.Background()

	// Table-driven test cases
	type args struct {
		userID int32
	}
	tests := []struct {
		name    string
		args    args
		want    *models.UserModel
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success - User Found",
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
			args: args{
				userID: 2,
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, name, is_delinquent FROM user WHERE id = ?`)).
					WithArgs(a.userID).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := repo.GetUserByID(ctx, tt.args.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}

			// Validate expectations
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}
