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

type paymentHandler struct {
	servicectx.ServiceCtx
}

func (p *paymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.PaymentRequest
	// Step 1: Decode the request body into PaymentRequest DTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorBadRequest).SetMessage("Request body is not valid").WriteResponse(w)
		return
	}

	// Step 2: Validate request (optional, can check if amount, user_id, loan_id, etc., are valid)
	if err := request.Validate(); err != nil {
		response.NewJSONResponse().SetError(errors.ErrorBadRequest).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	// Step 4: Create the payment record in the database
	err = p.ServiceCtx.PaymentService.CreatePayment(context.Background(), &request)
	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorInternalServer).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	// TODO push to rabbitMQ
	// Step 5: Respond with success message
	response.NewJSONResponse().SetData(nil).SetMessage("Payment successfully created").WriteResponse(w)
}

func NewPaymentHandler(ctx servicectx.ServiceCtx) PaymentHandlerInterface {
	return &paymentHandler{ctx}
}

type PaymentHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
}
