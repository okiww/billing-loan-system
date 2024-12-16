package servicectx

import (
	"github.com/okiww/billing-loan-system/internal/loan/services"
	userService "github.com/okiww/billing-loan-system/internal/user/services"
)

type ServiceCtx struct {
	LoanService services.LoanServiceInterface
	UserService userService.UserServiceInterface
}
