package handlerctx

import (
	"github.com/okiww/billing-loan-system/port/rest/handlers"
)

type HandlerCtx struct {
	LoanHandler handlers.LoanHandlerInterface
}
