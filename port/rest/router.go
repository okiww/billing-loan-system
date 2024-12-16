package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/okiww/billing-loan-system/port/rest/handlerctx"
)

type Domain struct {
	Domain handlerctx.HandlerCtx
}

// RegisterRoutes defines all application routes
func RegisterRoutes(router *mux.Router, h Domain) {
	baseRouter := router.PathPrefix("/api/v1").Subrouter()

	loanRouter := baseRouter.PathPrefix("/loan").Subrouter()
	loanRouter.HandleFunc("/create", h.Domain.LoanHandler.Create).Methods(http.MethodPost)

	paymentRouter := baseRouter.PathPrefix("/payment").Subrouter()
	paymentRouter.HandleFunc("/create", h.Domain.PaymentHandler.Create).Methods(http.MethodPost)
}
