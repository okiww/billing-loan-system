package rest

import (
	"github.com/gorilla/mux"
	"github.com/okiww/billing-loan-system/port/rest/handlers"
	"net/http"
)

// RegisterRoutes defines all application routes
func RegisterRoutes(router *mux.Router) {
	loanHandler := handlers.NewLoanHandler()
	loanRouter := router.PathPrefix("/loan").Subrouter()
	loanRouter.HandleFunc("/create", loanHandler.Create).Methods(http.MethodGet)
}
