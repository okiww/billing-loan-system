package handlerctx

import (
	"github.com/okiww/billing-loan-system/port/rest/handlers"
)

type HandlerCtx struct {
	LoanHandler    handlers.LoanHandlerInterface
	PaymentHandler handlers.PaymentHandlerInterface
}
