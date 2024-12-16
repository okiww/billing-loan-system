/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/okiww/billing-loan-system/configs"
	"github.com/okiww/billing-loan-system/internal/ctx/servicectx"
	"github.com/okiww/billing-loan-system/internal/loan/repositories"
	"github.com/okiww/billing-loan-system/internal/loan/services"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/robfig/cron/v3"
	"time"

	"github.com/spf13/cobra"
)

// backgroundCmd represents the background command
var backgroundCmd = &cobra.Command{
	Use:   "background",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runCronJob()
	},
}

func init() {
	rootCmd.AddCommand(backgroundCmd)
}

func runCronJob() {
	// Create a new cron scheduler
	c := cron.New(cron.WithLocation(time.Local))

	_, err := c.AddFunc("*/1 * * * *", GenerateBillPaymentEveryWeek) // Should change to run every monday at 00:00
	if err != nil {
		logger.GetLogger().Fatal("Error adding cron job:", err)
	}

	// Start the cron scheduler in the background
	c.Start()

	// Block main goroutine so that the program continues running
	// It will keep running indefinitely, and cron jobs will run as per the schedule
	select {}
}

// This function will be called by the cron job
func GenerateBillPaymentEveryWeek() {
	logger.GetLogger().Info("[Cronjob] Running cron for update bills")
	cfg := configs.InitConfig()
	// initial connection to database
	dbInit := mysql.InitDB(&cfg.DB)
	db, err := dbInit.Connect()
	if err != nil {
		logger.Fatalf("failed to connect db")
	}

	// initial domain context
	loanRepository := repositories.NewLoanRepository(db)
	loanBillRepository := repositories.NewLoanBillRepository(db)

	serviceCtx := servicectx.ServiceCtx{
		LoanService: services.NewLoanService(loanRepository, loanBillRepository),
	}

	ctx := context.Background()

	logger.GetLogger().Info("[Cronjob] Fetch all active loan")
	loans, err := serviceCtx.LoanService.GetAllActiveLoan(ctx)
	if err != nil {
		logger.Fatalf("[Cronjob] Error update loan bills")
		return
	}

	if len(loans) > 0 {
		// UpdateLoanBill update loan bill status
		logger.GetLogger().Info("[Cronjob] Update loan bill statuses")
		err = serviceCtx.LoanService.UpdateLoanBill(ctx)
		if err != nil {
			logger.Fatalf("[Cronjob] Error update loan bills")
			return
		}

		logger.GetLogger().Info("[Cronjob] Count loan bill overdue by loan id")
		for _, v := range loans {
			total, err := serviceCtx.LoanService.CountLoanBillOverdueStatusesByID(ctx, int32(v.ID))
			if err != nil {
				logger.Fatalf("[Cronjob] Error Count loan bill overdue by loan id")
				return
			}

			if total > 1 {
				// Update user is_delinquent true
			}
		}

	} else {
		logger.GetLogger().Info("[Cronjob] There's no active loan at the moment")
	}

	logger.GetLogger().Info("[Cronjob] Cronjob done")
	return
}
