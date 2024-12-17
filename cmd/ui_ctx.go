package cmd

import (
	"github.com/okiww/billing-loan-system/configs"
	billingConfigRepo "github.com/okiww/billing-loan-system/internal/billing_config/repositories"
	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	"github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/loan/services"
	paymentRepo "github.com/okiww/billing-loan-system/internal/payment/repositories"
	paymentService "github.com/okiww/billing-loan-system/internal/payment/services"
	userRepo "github.com/okiww/billing-loan-system/internal/user/repositories"
	userService "github.com/okiww/billing-loan-system/internal/user/services"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/mq"
	"github.com/okiww/billing-loan-system/port/rest/handlerctx"
	"github.com/okiww/billing-loan-system/port/rest/handlers"
)

func InitCtx(db *mysql.DBMySQL, mq *mq.RabbitMQ, rabbitMQCfg *configs.RabbitMQConfig) handlerctx.HandlerCtx {
	loanRepository := repositories.NewLoanRepository(db)
	loanBillRepository := repositories.NewLoanBillRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	paymentRepository := paymentRepo.NewPaymentRepository(db)
	billingConfigRepository := billingConfigRepo.NewBillingConfigRepository(db)

	serviceCtx := servicectx.ServiceCtx{
		LoanService:    services.NewLoanService(loanRepository, loanBillRepository, billingConfigRepository),
		UserService:    userService.NewUserService(userRepository),
		PaymentService: paymentService.NewPaymentService(paymentRepository, loanRepository, loanBillRepository),
	}

	handlerCtx := handlerctx.HandlerCtx{
		LoanHandler:    handlers.NewLoanHandler(serviceCtx),
		PaymentHandler: handlers.NewPaymentHandler(serviceCtx, mq, rabbitMQCfg),
	}

	return handlerCtx
}
