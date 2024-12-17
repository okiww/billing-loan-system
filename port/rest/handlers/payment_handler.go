package handlers

import (
	"context"
	"encoding/json"
	"github.com/okiww/billing-loan-system/configs"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/okiww/billing-loan-system/pkg/mq"
	"net/http"

	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	"github.com/okiww/billing-loan-system/internal/dto"
	"github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/okiww/billing-loan-system/pkg/response"
)

type paymentHandler struct {
	servicectx.ServiceCtx
	*mq.RabbitMQ
	*configs.RabbitMQConfig
}

func (p *paymentHandler) TestPublishMessage(w http.ResponseWriter, r *http.Request) {
	// Create an array of Payment structs
	payments := []models.Payment{
		{ID: 1, UserID: 101, LoanID: 202, LoanBillID: 303, Amount: 5000, TotalAmount: 5100, Status: "paid"},
		{ID: 2, UserID: 102, LoanID: 203, LoanBillID: 304, Amount: 3000, TotalAmount: 3100, Status: "pending"},
		{ID: 3, UserID: 103, LoanID: 204, LoanBillID: 305, Amount: 7000, TotalAmount: 7100, Status: "paid"},
	}

	// Serialize the array to JSON
	jsonData, err := json.Marshal(payments)
	if err != nil {
		logger.GetLogger().Fatalf("Failed to marshal array to JSON: %v", err)
	}

	err = p.RabbitMQ.PublishMessage(p.RabbitMQConfig.QueueName, string(jsonData))
	if err != nil {
		logger.GetLogger().Fatalf("Failed to publish message: %v", err)
	}

	logger.GetLogger().Println("Message published successfully!")
	response.NewJSONResponse().SetData(nil).SetMessage("Message published successfully!").WriteResponse(w)
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
	err = p.ServiceCtx.PaymentService.MakePayment(context.Background(), &request)
	if err != nil {
		response.NewJSONResponse().SetError(errors.ErrorInternalServer).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	// TODO push to rabbitMQ
	// Step 5: Respond with success message
	response.NewJSONResponse().SetData(nil).SetMessage("Payment successfully created").WriteResponse(w)
}

func NewPaymentHandler(ctx servicectx.ServiceCtx, rabbitMQ *mq.RabbitMQ, rabbitMQCfg *configs.RabbitMQConfig) PaymentHandlerInterface {
	return &paymentHandler{ctx, rabbitMQ, rabbitMQCfg}
}

type PaymentHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	TestPublishMessage(w http.ResponseWriter, r *http.Request)
}
