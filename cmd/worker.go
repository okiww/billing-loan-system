/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"github.com/okiww/billing-loan-system/configs"
	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	loanRepo "github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/payment/models"
	paymentRepo "github.com/okiww/billing-loan-system/internal/payment/repositories"
	"github.com/okiww/billing-loan-system/internal/payment/services"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/okiww/billing-loan-system/pkg/mq"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitWorker()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

func InitWorker() {
	cfg := configs.InitConfig()

	// initial connection to database
	dbInit := mysql.InitDB(&cfg.DB)
	db, err := dbInit.Connect()
	if err != nil {
		logger.Fatalf("failed to connect db")
	}

	// initial connection to rabbitMQ
	rabbitMQ, err := mq.NewRabbitMQ(cfg.RabbitMQ.Dsn)
	if err != nil {
		logger.GetLogger().Fatalf("failed to connect to RabbitMQ: %v", err)
		return
	}
	defer rabbitMQ.Close() // Ensure connection is closed when the function exits

	_, err = rabbitMQ.DeclareQueue(cfg.RabbitMQ.QueueName)
	if err != nil {
		logger.GetLogger().Fatalf("failed to declare queue %s: %v", cfg.RabbitMQ.QueueName, err)
		return
	}

	// initial domain context
	loanRepository := loanRepo.NewLoanRepository(db)
	paymentRepository := paymentRepo.NewPaymentRepository(db)

	serviceCtx := servicectx.ServiceCtx{
		PaymentService: services.NewPaymentService(paymentRepository, loanRepository),
	}

	messages, err := rabbitMQ.ConsumeMessages(cfg.RabbitMQ.QueueName)
	if err != nil {
		logger.GetLogger().Fatalf("failed to consume messages from queue %s: %v", cfg.RabbitMQ.QueueName, err)
		return
	}

	logger.GetLogger().Info("Worker started, waiting for messages. Press CTRL+C to stop.")

	// Channel to listen for OS signals (e.g., SIGINT, SIGTERM)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Process messages in a goroutine
	go func() {
		for msg := range messages {
			logger.GetLogger().Infof("Received message: %s", string(msg.Body))
			err := processPayment(context.Background(), &serviceCtx, msg.Body)
			if err != nil {
				logger.GetLogger().Errorf("Failed to process message: %v", err)
			} else {
				logger.GetLogger().Info("Message processed successfully")
			}
		}
	}()

	// Wait for termination signal
	<-signalChan
	logger.GetLogger().Info("Graceful shutdown: worker stopping...")
}

// processPayment processes the incoming RabbitMQ message body for payment
func processPayment(ctx context.Context, serviceCtx *servicectx.ServiceCtx, body []byte) error {
	var payments models.Payment
	err := json.Unmarshal(body, &payments)
	if err != nil {
		logger.GetLogger().Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	err = serviceCtx.PaymentService.ProcessUpdatePayment(ctx, payments)
	if err != nil {
		logger.GetLogger().Errorf("Failed to process payment: %v", err)
		return err
	}

	// Log the received array of Payment structs
	logger.GetLogger().Infof("Done Process Payment: %+v", payments)
	return nil
}
