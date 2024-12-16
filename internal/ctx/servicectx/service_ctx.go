package servicectx

import (
	"github.com/okiww/billing-loan-system/internal/loan/services"
	services2 "github.com/okiww/billing-loan-system/internal/payment/services"
	userService "github.com/okiww/billing-loan-system/internal/user/services"
)

type ServiceCtx struct {
	LoanService    services.LoanServiceInterface
	UserService    userService.UserServiceInterface
	PaymentService services2.PaymentServiceInterface
}
