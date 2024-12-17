/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/okiww/billing-loan-system/pkg/mq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/okiww/billing-loan-system/configs"
	mysql "github.com/okiww/billing-loan-system/pkg/db"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/okiww/billing-loan-system/port/rest"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ServeHttp()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func ServeHttp() {
	cfg := configs.InitConfig()
	// initial connection to database
	dbInit := mysql.InitDB(&cfg.DB)
	db, err := dbInit.Connect()
	if err != nil {
		log.Fatalf("failed to connect db")
	}

	// initial rabbitMQ
	rabbitMQ, err := mq.NewRabbitMQ(cfg.RabbitMQ.Dsn)
	if err != nil {
		logger.GetLogger().Fatalf("failed to connect to RabbitMQ: %v", err)
		return
	}
	defer rabbitMQ.Close()

	// initial domain context
	domainCtx := InitCtx(db, rabbitMQ, &cfg.RabbitMQ)

	// initial router
	router := mux.NewRouter()
	rest.RegisterRoutes(router, rest.Domain{
		Domain: domainCtx,
	})

	// Create the HTTP server
	server := &http.Server{
		Addr:    cfg.Http.Addr,
		Handler: router,
	}

	// Run the server in a separate goroutine
	go func() {
		log.Println("Server running on port", cfg.Http.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Set up channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until a signal is received
	<-stop
	log.Println("Shutting down server...")

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	defer func() {
		// extra handling here
		err := db.CloseDB()
		if err != nil {
			logger.Fatalf("failed close db %s", err.Error())
		}
		serverStopCtx()
		<-serverCtx.Done()
	}()

	// Attempt a graceful shutdown
	if err := server.Shutdown(serverCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
