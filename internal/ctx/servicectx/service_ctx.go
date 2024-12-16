package servicectx

import (
	"github.com/okiww/billing-loan-system/internal/loan/services"
)

type ServiceCtx struct {
	LoanService services.LoanServiceInterface
}
