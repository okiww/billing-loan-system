package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/okiww/billing-loan-system/helpers"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
)

var (
	repo     PaymentRepositoryInterface
	repoLock sync.Once
)

type paymentRepository struct {
	*mysql.DBMySQL
}

func (p *paymentRepository) Create(ctx context.Context, payment *models.Payment) (int32, error) {
	query := `
		INSERT INTO payments (user_id, loan_id, loan_bill_id, amount, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := p.DB.ExecContext(ctx, query, payment.UserID, payment.LoanID, payment.LoanBillID, payment.Amount, payment.Status, time.Now())
	if err != nil {
		return 0, err
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

func (p *paymentRepository) UpdatePaymentStatus(ctx context.Context, id int32, status string, note string) error {
	query := `
		UPDATE payments SET status = ?, updated_at = ?, note = ? WHERE id = ?
	`
	_, err := p.DB.ExecContext(ctx, query, status, time.Now(), note, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *paymentRepository) GetPaymentByID(ctx context.Context, id int32) (*models.Payment, error) {

	query := "SELECT id, user_id, loan_id, loan_bill_id, amount, status, created_at, updated_at, note FROM payments WHERE id = ?"

	rows, err := p.DB.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payment models.Payment
	for rows.Next() {
		if err := rows.StructScan(&payment); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &payment, nil
}

type PaymentRepositoryInterface interface {
	Create(ctx context.Context, payment *models.Payment) (int32, error)
	UpdatePaymentStatus(ctx context.Context, id int32, status string, note string) error
	GetPaymentByID(ctx context.Context, id int32) (*models.Payment, error)
}

func NewPaymentRepository(db *mysql.DBMySQL) PaymentRepositoryInterface {
	if helpers.IsTestEnv() { // Skip singleton in tests
		return &paymentRepository{
			db,
		}
	}

	repoLock.Do(func() {
		repo = &paymentRepository{
			db,
		}
	})
	return repo
}
