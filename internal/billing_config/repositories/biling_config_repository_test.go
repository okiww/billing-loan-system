package repositories

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/okiww/billing-loan-system/internal/billing_config/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestBillingConfigRepository_GetBillingConfigByName(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	mockDB := &mysql.DBMySQL{DB: db}
	repo := NewBillingConfigRepository(mockDB)

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		s       BillingConfigRepositoryInterface
		args    args
		want    *models.BillingConfig
		wantErr bool
		mock    func(a args)
	}{
		{
			name: "Success",
			s:    repo,
			args: args{
				name: "loan_interest_percentage",
			},
			want: &models.BillingConfig{
				ID:    1,
				Name:  "loan_interest_percentage",
				Value: "{\"is_active\":true,\"value\":10}",
			},
			wantErr: false,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM billing_configs WHERE name = ?`)).
					WithArgs(a.name).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "value",
					}).AddRow(
						1, "loan_interest_percentage", "{\"is_active\":true,\"value\":10}",
					))
			},
		},
		{
			name: "Config Not Found",
			s:    repo,
			args: args{
				name: "loan_interest_percentage",
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM billing_configs WHERE name = ?`)).
					WithArgs(a.name).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Database Error",
			s:    repo,
			args: args{
				name: "loan_interest_percentage",
			},
			want:    nil,
			wantErr: true,
			mock: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM billing_configs WHERE name = ?`)).
					WithArgs(a.name).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)

			got, err := tt.s.GetBillingConfigByName(context.Background(), tt.args.name)

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
