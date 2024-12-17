package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/okiww/billing-loan-system/configs"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/okiww/billing-loan-system/pkg/mq"

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
	payments := models.Payment{
		ID: 3, UserID: 123, LoanID: 2, LoanBillID: 1, Amount: 1375000, Status: "PENDING",
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
	payment, err := p.ServiceCtx.PaymentService.MakePayment(context.Background(), &request)
	if err != nil {
		if err.Error() == dto.ErrorLoanIsNotActive || err.Error() == dto.ErrorPaymentAmountNotMatchWithBill {
			response.NewJSONResponse().SetError(errors.ErrorBadRequest).SetMessage(err.Error()).WriteResponse(w)
			return
		}
		response.NewJSONResponse().SetError(errors.ErrorInternalServer).SetMessage(err.Error()).WriteResponse(w)
		return
	}

	// Step 5: Push to rabbitMQ
	go p.publishPayment(payment)
	// Step 6: Respond with success message
	response.NewJSONResponse().SetData(nil).SetMessage("Payment successfully created").WriteResponse(w)
}

func (p *paymentHandler) publishPayment(payment *models.Payment) {
	// Serialize the array to JSON
	jsonData, err := json.Marshal(payment)
	if err != nil {
		logger.GetLogger().Fatalf("Failed to marshal array to JSON: %v", err)
	}

	err = p.RabbitMQ.PublishMessage(p.RabbitMQConfig.QueueName, string(jsonData))
	if err != nil {
		logger.GetLogger().Fatalf("Failed to publish message: %v", err)
	}

	logger.GetLogger().Println("Message published successfully!")
}

func NewPaymentHandler(ctx servicectx.ServiceCtx, rabbitMQ *mq.RabbitMQ, rabbitMQCfg *configs.RabbitMQConfig) PaymentHandlerInterface {
	return &paymentHandler{ctx, rabbitMQ, rabbitMQCfg}
}

type PaymentHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	TestPublishMessage(w http.ResponseWriter, r *http.Request)
}
