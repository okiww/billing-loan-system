package cmd

import (
	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	"github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/loan/services"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/port/rest/handlerctx"
	"github.com/okiww/billing-loan-system/port/rest/handlers"
)

func InitCtx(db *mysql.DBMySQL) handlerctx.HandlerCtx {
	loanRepository := repositories.NewLoanRepository(db)

	serviceCtx := servicectx.ServiceCtx{
		LoanService: services.NewLoanService(loanRepository),
	}

	handlerCtx := handlerctx.HandlerCtx{
		LoanHandler: handlers.NewLoanHandler(serviceCtx),
	}

	return handlerCtx
}
