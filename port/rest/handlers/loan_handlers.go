package handlers

import (
	"context"
	"encoding/json"
	"net/http"

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

func NewLoanHandler(ctx servicectx.ServiceCtx) LoanHandlerInterface {
	return &loanHandler{ctx}
}

type LoanHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
}
