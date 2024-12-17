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

func (p *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (user_id, loan_id, loan_bill_id, amount, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := p.DB.ExecContext(ctx, query, payment.UserID, payment.LoanID, payment.LoanBillID, payment.Amount, payment.Status, time.Now())
	if err != nil {
		return err
	}
	return nil
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

type PaymentRepositoryInterface interface {
	Create(ctx context.Context, payment *models.Payment) error
	UpdatePaymentStatus(ctx context.Context, id int32, status string, note string) error
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