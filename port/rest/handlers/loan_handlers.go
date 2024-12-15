package handlers

import (
	"fmt"
	"net/http"
)

type loanHandler struct{}

func (l loanHandler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apaci apaci")
}

func NewLoanHandler() LoanHandlerInterface {
	return &loanHandler{}
}

type LoanHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
}
