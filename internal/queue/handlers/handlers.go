package handlers

import (
	"context"
	"encoding/json"
	"github.com/Uranury/RBK_finalProject/internal/queue/jobs"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/hibiken/asynq"
	"log/slog"
)

type WorkerHandler struct {
	EmailService   *services.EmailService
	InvoiceService *services.InvoiceService
	asynqClient    *asynq.Client
	logger         *slog.Logger
}

func NewWorkerHandler(emailService *services.EmailService, invoiceService *services.InvoiceService, logger *slog.Logger) *WorkerHandler {
	return &WorkerHandler{EmailService: emailService, InvoiceService: invoiceService, logger: logger}
}

func (h *WorkerHandler) HandleGeneratePDFTask(ctx context.Context, t *asynq.Task) error {
	var payload jobs.PDFPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		h.logger.Warn("failed to unmarshal PDF payload", "err", err)
		return err
	}

	pdfBytes, err := h.InvoiceService.GenerateInvoicePDF(ctx, payload.OrderID, payload.OrderItemID)
	if err != nil {
		h.logger.Warn("failed to generate PDF", "err", err)
		return err
	}

	emailTask, err := jobs.NewSendEmailTask("recipient@example.com", "Your Invoice", "Please find attached.", pdfBytes)
	if err != nil {
		h.logger.Error("failed to create send-email task", "err", err)
		return err
	}

	if _, err := h.asynqClient.Enqueue(emailTask); err != nil {
		h.logger.Error("failed to enqueue email task", "err", err)
		return err
	}

	return nil
}
