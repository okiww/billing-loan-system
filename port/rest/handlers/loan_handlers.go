package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	"github.com/okiww/billing-loan-system/internal/dto"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/okiww/billing-loan-system/pkg/response"
)

type loanHandler struct {
	servicectx.ServiceCtx
}

func (l *loanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.LoanRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorBadRequest).SetMessage("Request body is not valid").WriteResponse(w)
		return
	}

	// validate request
	if err := request.Validate(); err != nil {
		response.NewJSONResponse().SetError(errors.ErrorBadRequest).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	// validate users isDelinquent or not
	isDelinquent, err := l.ServiceCtx.UserService.IsDelinquent(context.Background(), int32(request.UserID))
	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorInternalServer).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	if isDelinquent {
		response.NewJSONResponse().SetError(errors.ErrorForbiddenResource).SetMessage("Couldn't create loan because user is delinquent").WriteResponse(w)
		return
	}

	err = l.ServiceCtx.LoanService.CreateLoan(context.Background(), request)
	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorInternalServer).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	response.NewJSONResponse().SetData(nil).SetMessage("Success create transaction").WriteResponse(w)
}

func (l *loanHandler) GetLoans(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorBadRequest).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	loansWithBills, err := l.LoanService.GetLoansWithBills(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch loans with bills", http.StatusInternalServerError)
		return
	}

	response.NewJSONResponse().SetData(loansWithBills).SetMessage("Success get loans").WriteResponse(w)
}

func NewLoanHandler(ctx servicectx.ServiceCtx) LoanHandlerInterface {
	return &loanHandler{ctx}
}

type LoanHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetLoans(w http.ResponseWriter, r *http.Request)
}
