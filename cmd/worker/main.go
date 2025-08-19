package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Uranury/RBK_finalProject/internal/queue/handlers"
	"github.com/Uranury/RBK_finalProject/internal/queue/jobs"
	"github.com/Uranury/RBK_finalProject/internal/repositories/order"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/hibiken/asynq"
	"github.com/mailgun/mailgun-go/v4"
)

func main() {
	// Use a more verbose logger that shows all levels
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	deps, err := InitWorkerDeps(logger)
	if err != nil {
		logger.Error("failed to init worker", "err", err)
		os.Exit(1)
	}

	mux := asynq.NewServeMux()

	// Initialize services used by worker handlers
	ordRepo := order.NewRepository(deps.DB)
	invoiceService := services.NewInvoiceService(ordRepo, deps.Logger)

	mg := mailgun.NewMailgun(deps.Cfg.MailgunDomain, deps.Cfg.MailgunAPIKey)
	emailService := services.NewEmailService(mg, deps.Cfg.MailgunDomain, deps.Logger)

	workerHandler := handlers.NewWorkerHandler(emailService, invoiceService, deps.Logger)

	mux.HandleFunc(jobs.SendInvoice, func(ctx context.Context, t *asynq.Task) error {
		logger.Info("processing send-invoice task", "task_id", t.ResultWriter().TaskID())
		return workerHandler.HandleSendInvoiceTask(ctx, t)
	})

	if err := deps.Server.Run(mux); err != nil {
		logger.Error("could not run asynq server", "err", err)
		os.Exit(1)
	}
}
