package cmd

import (
	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	"github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/loan/services"
	userRepo "github.com/okiww/billing-loan-system/internal/user/repositories"
	userService "github.com/okiww/billing-loan-system/internal/user/services"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/port/rest/handlerctx"
	"github.com/okiww/billing-loan-system/port/rest/handlers"
)

func InitCtx(db *mysql.DBMySQL) handlerctx.HandlerCtx {
	loanRepository := repositories.NewLoanRepository(db)
	loanBillRepository := repositories.NewLoanBillRepository(db)
	userRepository := userRepo.NewUserRepository(db)

	serviceCtx := servicectx.ServiceCtx{
		LoanService: services.NewLoanService(loanRepository, loanBillRepository),
		UserService: userService.NewUserService(userRepository),
	}

	handlerCtx := handlerctx.HandlerCtx{
		LoanHandler: handlers.NewLoanHandler(serviceCtx),
	}

	return handlerCtx
}
